package main

import (
	"errors"
	"sort"
	"time"
)

func getLastVisit(visits map[string][]visit, cmd command) (time.Time, error) {
	patientName := cmd.data[0]
	doctor := cmd.data[1]

	// retrieve all visits for the patient across all doctors
	history := visits[patientName]
	if len(history) == 0 {
		return time.Time{}, patientNotFoundError{}
	}

	var doctorVisits []visit
	for _, v := range history {
		if v.doctor == doctor {
			doctorVisits = append(doctorVisits, v)
		}
	}
	if len(doctorVisits) == 0 { // no visits to this specific doctor
		return time.Time{}, errors.New("no visits found for doctor")
	}

	// sort visits by date in ascending order to find the most recent
	sort.Slice(doctorVisits, func(i, j int) bool {
		return doctorVisits[i].date.Before(doctorVisits[j].date)
	})

	// return the last element (most recent date) after sorting
	return doctorVisits[len(doctorVisits)-1].date, nil
}

func getHistory(visits map[string][]visit, cmd command) ([]visit, error) {
	history := visits[cmd.data[0]]
	if len(history) == 0 {
		return nil, patientNotFoundError{} // patient has no recorded visits
	}

	// sort visits by date in chronological order
	sort.Slice(history, func(i, j int) bool {
		return history[i].date.Before(history[j].date)
	})

	return history, nil
}

func setVisit(visits map[string][]visit, cmd command) {
	patientName := cmd.data[0]
	doctor := cmd.data[1]
	date, _ := time.Parse("2006-01-02", cmd.data[2])
	newVisit := visit{
		doctor: doctor,
		date:   date,
	}
	visits[patientName] = append(visits[patientName], newVisit)
}
