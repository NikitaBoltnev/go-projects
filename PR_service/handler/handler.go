package handler

import (
	"PR_service/api"
	"context"
	"encoding/json"
	"math/rand"
	"slices"

	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(dbPool *pgxpool.Pool) *Handler {
	return &Handler{
		db: dbPool,
	}
}

func (h *Handler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var reqBody api.PostPullRequestCreateJSONBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback(r.Context())
	}()

	sql1 := "SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1 OR pull_request_name = $2)"
	var exists bool
	err = tx.QueryRow(r.Context(), sql1, reqBody.PullRequestId, reqBody.PullRequestName).Scan(&exists)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if exists {
		sendErrorResponse(w, api.PREXISTS, "PR id already exists", http.StatusConflict)
		return
	}

	authorTeamName, err := h.getUserTeamName(tx, r.Context(), reqBody.AuthorId)
	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "author not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	_, err = h.getTeamMembers(tx, r.Context(), authorTeamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "author team not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	passedUserIds := []string{reqBody.AuthorId}
	activeRev, err := h.getActiveReviewersFromTeam(tx, r.Context(), authorTeamName, passedUserIds)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	shuffleStr(activeRev)
	if len(activeRev) > 2 {
		activeRev = activeRev[:2]
	}

	finalRevIds := activeRev

	if finalRevIds == nil {
		finalRevIds = []string{}
	}

	prInsertSQL := `
    INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, assigned_reviewer_ids, created_at) 
    VALUES ($1, $2, $3, $4, $5, $6);`
	now := time.Now()

	_, err = tx.Exec(r.Context(), prInsertSQL,
		reqBody.PullRequestId,
		reqBody.PullRequestName,
		reqBody.AuthorId,
		api.PullRequestStatusOPEN,
		finalRevIds,
		now,
	)

	if err != nil {
		http.Error(w, "Failed to create PR", http.StatusInternalServerError)
		return
	}

	createdPR := api.PullRequest{
		PullRequestId:     reqBody.PullRequestId,
		PullRequestName:   reqBody.PullRequestName,
		AuthorId:          reqBody.AuthorId,
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: finalRevIds,
		CreatedAt:         &now,
		MergedAt:          nil,
	}

	responseBody := map[string]any{
		"pr": createdPR,
	}

	err = tx.Commit(r.Context())
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, responseBody, http.StatusCreated)
}

func (h *Handler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var reqBody api.PostPullRequestMergeJSONRequestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback(r.Context())
	}()

	var currentStatus api.PullRequestStatus
	var currentPRData api.PullRequest

	err = tx.QueryRow(r.Context(), `
        SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewer_ids, created_at, merged_at
        FROM pull_requests
        WHERE pull_request_id = $1`,
		reqBody.PullRequestId).Scan(
		&currentPRData.PullRequestId,
		&currentPRData.PullRequestName,
		&currentPRData.AuthorId,
		&currentStatus,
		&currentPRData.AssignedReviewers,
		&currentPRData.CreatedAt,
		&currentPRData.MergedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "PR not found", http.StatusNotFound)
			return
		}

		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	var mergedAtToReturn *time.Time

	if currentStatus != api.PullRequestStatusMERGED {
		_, err = tx.Exec(r.Context(), `
            UPDATE pull_requests
            SET status = $1, merged_at = $2
            WHERE pull_request_id = $3`,
			api.PullRequestStatusMERGED,
			&now,
			reqBody.PullRequestId,
		)
		if err != nil {

			http.Error(w, "failed to merge PR", http.StatusInternalServerError)
			return
		}
		mergedAtToReturn = &now
	} else {
		mergedAtToReturn = currentPRData.MergedAt
	}

	responsePR := api.PullRequest{
		PullRequestId:     currentPRData.PullRequestId,
		PullRequestName:   currentPRData.PullRequestName,
		AuthorId:          currentPRData.AuthorId,
		Status:            api.PullRequestStatusMERGED,
		AssignedReviewers: currentPRData.AssignedReviewers,
		CreatedAt:         currentPRData.CreatedAt,
		MergedAt:          mergedAtToReturn,
	}

	responseBody := map[string]any{
		"pr": responsePR,
	}

	err = tx.Commit(r.Context())
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, responseBody, http.StatusOK)
}

func (h *Handler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var reqBody api.PostPullRequestReassignJSONBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback(r.Context())
	}()

	var currentPRData api.PullRequest
	var currentStatus api.PullRequestStatus
	var assignedReviewerIds []string

	err = tx.QueryRow(r.Context(), `
        SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewer_ids, created_at, merged_at
        FROM pull_requests
        WHERE pull_request_id = $1`,
		reqBody.PullRequestId).Scan(
		&currentPRData.PullRequestId,
		&currentPRData.PullRequestName,
		&currentPRData.AuthorId,
		&currentStatus,
		&assignedReviewerIds,
		&currentPRData.CreatedAt,
		&currentPRData.MergedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "PR not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if currentStatus == api.PullRequestStatusMERGED {
		sendErrorResponse(w, api.PRMERGED, "cannot reassign on merged PR", http.StatusConflict)
		return
	}

	oldUserIndex, found := findReviewerIndex(assignedReviewerIds, reqBody.OldUserId)
	if !found {
		sendErrorResponse(w, api.NOTASSIGNED, "reviewer is not assigned to this PR", http.StatusConflict)
		return
	}

	var oldUserTeamName string
	oldUserTeamName, err = h.getUserTeamName(tx, r.Context(), reqBody.OldUserId)
	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "reviewer user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	excludeUserIds := h.buildReviewerExclusionListForReassign(reqBody.OldUserId, currentPRData.AuthorId, assignedReviewerIds)

	activeCandidates, err := h.getActiveReviewersFromTeam(tx, r.Context(), oldUserTeamName, excludeUserIds)
	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "reviewer team not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if len(activeCandidates) == 0 {
		sendErrorResponse(w, api.NOCANDIDATE, "no active replacement candidate in team", http.StatusConflict)
		return
	}

	shuffleStr(activeCandidates)
	newUser := activeCandidates[0]

	newReviewerIds := make([]string, 0, len(assignedReviewerIds))
	for i, id := range assignedReviewerIds {
		if i != oldUserIndex {
			newReviewerIds = append(newReviewerIds, id)
		}
	}
	newReviewerIds = append(newReviewerIds, newUser)

	_, err = tx.Exec(r.Context(), `
        UPDATE pull_requests
        SET assigned_reviewer_ids = $1
        WHERE pull_request_id = $2`,
		newReviewerIds,
		reqBody.PullRequestId,
	)

	if err != nil {
		http.Error(w, "failed to reassign reviewer", http.StatusInternalServerError)
		return
	}

	currentPRData.AssignedReviewers = newReviewerIds

	responseBody := map[string]any{
		"pr":          currentPRData,
		"replaced_by": newUser,
	}

	err = tx.Commit(r.Context())
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, responseBody, http.StatusOK)
}

func (h *Handler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var team api.Team
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&team)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback(r.Context())
	}()

	var existingTeamName string
	err = tx.QueryRow(r.Context(), "SELECT team_name FROM teams WHERE team_name = $1", team.TeamName).Scan(&existingTeamName)
	if err == nil {
		sendErrorResponse(w, api.TEAMEXISTS, "team_name already exists", http.StatusBadRequest)
		return
	} else if err != pgx.ErrNoRows {

		http.Error(w, "failed to check team existence", http.StatusInternalServerError)
		return
	}

	memberIds := make([]string, len(team.Members))
	for i, member := range team.Members {
		memberIds[i] = member.UserId
	}

	usersOldTeams, err := h.getUsersOldTeamsFromTeamsTable(r.Context(), memberIds)
	if err != nil {

		http.Error(w, "failed to check old team memberships (teams)", http.StatusInternalServerError)
		return
	}

	err = h.upsertUsers(r.Context(), team.Members, team.TeamName)
	if err != nil {
		http.Error(w, "failed to add/update users", http.StatusInternalServerError)
		return
	}

	err = h.removeUsersFromOldTeams(r.Context(), usersOldTeams, team.TeamName)
	if err != nil {
		http.Error(w, "failed to update old teams", http.StatusInternalServerError)
		return
	}

	sql := "INSERT INTO teams (team_name, members) VALUES ($1, $2)"

	_, err = tx.Exec(r.Context(), sql, team.TeamName, memberIds)
	if err != nil {
		http.Error(w, "failed to add team", http.StatusInternalServerError)
		return
	}

	responseBody := map[string]any{
		"team": team,
	}

	err = tx.Commit(r.Context())
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, responseBody, http.StatusCreated)
}

func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	tx, err := h.db.Begin(r.Context())
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback(r.Context())
	}()

	userIds, err := h.getTeamMembers(tx, r.Context(), params.TeamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "team_name not found", http.StatusNotFound)
			return
		}

		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if len(userIds) == 0 {
		resultTeam := api.Team{
			TeamName: params.TeamName,
			Members:  []api.TeamMember{},
		}
		sendJSONResponse(w, resultTeam, http.StatusOK)

		return
	}

	sql := "SELECT  user_id, username, is_active FROM users WHERE user_id = ANY($1)"
	rows, err := tx.Query(r.Context(), sql, userIds)
	if err != nil {

		http.Error(w, "failed to get team members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var members []api.TeamMember
	for rows.Next() {
		var user api.User
		err = rows.Scan(&user.UserId, &user.Username, &user.IsActive)
		if err != nil {

			http.Error(w, "failed to process team members", http.StatusInternalServerError)
			return
		}
		members = append(members, api.TeamMember{
			UserId:   user.UserId,
			Username: user.Username,
			IsActive: user.IsActive,
		})

	}

	if err = rows.Err(); err != nil {
		http.Error(w, "failed to process team members", http.StatusInternalServerError)
		return
	}
	if members == nil {
		members = []api.TeamMember{}
	}
	resultTeam := api.Team{
		TeamName: params.TeamName,
		Members:  members,
	}

	err = tx.Commit(r.Context())
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, resultTeam, http.StatusOK)
}

func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	userId := params.UserId
	rows, err := h.db.Query(r.Context(), `
		SELECT pull_request_id, pull_request_name, author_id, status
		FROM pull_requests
		WHERE $1 = ANY(assigned_reviewer_ids);`,
		userId)

	if err != nil {

		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pullRequests []api.PullRequestShort
	for rows.Next() {
		var prShort api.PullRequestShort
		err = rows.Scan(&prShort.PullRequestId, &prShort.PullRequestName, &prShort.AuthorId, &prShort.Status)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		pullRequests = append(pullRequests, prShort)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if pullRequests == nil {
		pullRequests = []api.PullRequestShort{}
	}
	responseBody := map[string]any{
		"user_id":       userId,
		"pull_requests": pullRequests,
	}

	sendJSONResponse(w, responseBody, http.StatusOK)

}

func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var reqBody api.PostUsersSetIsActiveJSONBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	sql := `
		UPDATE users
		SET is_active = $2
		WHERE user_id = $1
		RETURNING user_id, username, team_name, is_active;
	`
	var updatedUser api.User
	err = h.db.QueryRow(r.Context(), sql, reqBody.UserId, reqBody.IsActive).Scan(
		&updatedUser.UserId,
		&updatedUser.Username,
		&updatedUser.TeamName,
		&updatedUser.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			sendErrorResponse(w, api.NOTFOUND, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	responseBody := map[string]any{
		"user": updatedUser,
	}

	sendJSONResponse(w, responseBody, http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, code api.ErrorResponseErrorCode, message string, statusCode int) {
	errorRes := api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorRes)
}

func sendJSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func shuffleStr(slice []string) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func (h *Handler) getActiveReviewersFromTeam(tx pgx.Tx, ctx context.Context, teamName string, excludeUserIds []string) ([]string, error) {
	members, err := h.getTeamMembers(tx, ctx, teamName)
	if err != nil {
		return nil, err
	}

	candidateIds := make([]string, 0, len(members))
	for _, memberId := range members {
		if !slices.Contains(excludeUserIds, memberId) {
			candidateIds = append(candidateIds, memberId)
		}
	}

	if len(candidateIds) == 0 {
		return []string{}, nil
	}

	rows, err := tx.Query(ctx, "SELECT user_id FROM users WHERE user_id = ANY($1) AND is_active = true", candidateIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeRev []string
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		activeRev = append(activeRev, userId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return activeRev, nil
}

func (h *Handler) getTeamMembers(tx pgx.Tx, ctx context.Context, teamName string) ([]string, error) {
	var members []string
	err := tx.QueryRow(ctx, "SELECT members FROM teams WHERE team_name = $1", teamName).Scan(&members)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return members, nil
}

func (h *Handler) getUserTeamName(tx pgx.Tx, ctx context.Context, userID string) (string, error) {
	var teamName string
	err := tx.QueryRow(ctx, "SELECT team_name FROM users WHERE user_id = $1", userID).Scan(&teamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", err
		}
		return "", err
	}
	return teamName, nil
}

func (h *Handler) buildReviewerExclusionListForReassign(oldUserId, authorId string, assignedReviewerIds []string) []string {
	passed := make(map[string]bool)
	passed[oldUserId] = true
	passed[authorId] = true

	for _, assignedId := range assignedReviewerIds {
		if assignedId != oldUserId {
			passed[assignedId] = true
		}
	}

	excludeList := make([]string, 0, len(passed))
	for id := range passed {
		excludeList = append(excludeList, id)
	}

	return excludeList
}

func findReviewerIndex(reviewers []string, userId string) (int, bool) {
	for i, id := range reviewers {
		if id == userId {
			return i, true
		}
	}
	return 0, false
}

func (h *Handler) getUsersOldTeamsFromTeamsTable(ctx context.Context, userIds []string) ([]struct{ UserId, OldTeam string }, error) {
	var mappings []struct{ UserId, OldTeam string }

	if len(userIds) == 0 {
		return mappings, nil
	}

	rows, err := h.db.Query(ctx, `
        SELECT t.team_name, m AS user_id_in_team
        FROM teams t, unnest(t.members) AS m
        WHERE m = ANY($1);
    `, userIds)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mapping struct {
			UserId  string
			OldTeam string
		}
		err = rows.Scan(&mapping.OldTeam, &mapping.UserId)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mappings, nil
}

func (h *Handler) removeUsersFromOldTeams(ctx context.Context, oldTeamMappings []struct{ UserId, OldTeam string }, newTeamName string) error {
	for _, mapping := range oldTeamMappings {
		if mapping.OldTeam != newTeamName {
			_, err := h.db.Exec(ctx, "UPDATE teams SET members = array_remove(members, $1) WHERE team_name = $2", mapping.UserId, mapping.OldTeam)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Handler) upsertUsers(ctx context.Context, members []api.TeamMember, teamName string) error {
	sql := `
        INSERT INTO users (user_id, username, is_active, team_name)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id)
        DO UPDATE SET
            is_active = EXCLUDED.is_active,
            team_name = EXCLUDED.team_name;`

	for _, member := range members {
		_, err := h.db.Exec(ctx, sql, member.UserId, member.Username, member.IsActive, teamName)
		if err != nil {
			return err
		}
	}
	return nil
}
