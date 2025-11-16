CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(100) PRIMARY KEY,
    username VARCHAR(100) NOT NULL, 
   team_name VARCHAR(100),     
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR(100) UNIQUE NOT NULL, 
    members VARCHAR(100)[] NOT NULL 
);

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(100) PRIMARY KEY,
    pull_request_name VARCHAR(200) NOT NULL,
    author_id VARCHAR(100) NOT NULL,
    status VARCHAR(10) DEFAULT 'OPEN',
    assigned_reviewer_ids VARCHAR(100)[],
    created_at TIMESTAMP NULL,
    merged_at TIMESTAMP NULL  
);

CREATE INDEX IF NOT EXISTS idx_teams_members_gin ON teams USING GIN (members);
CREATE INDEX IF NOT EXISTS idx_pr_author_id ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_assigned_reviewers_gin ON pull_requests USING GIN (assigned_reviewer_ids);
