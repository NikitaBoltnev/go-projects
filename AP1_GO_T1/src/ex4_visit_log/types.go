package main

import "time"

// represents an error when a patient record cannot be found in the system
type patientNotFoundError struct {
}

// error interface for patientNotFoundError
func (e patientNotFoundError) Error() string {
	return "patient not found"
}

// represents a user command with associated data and error state
// used for communication between input processing and command execution layers
// name: command type identifier ("save", "getHistory", "getLastVisit")
// data: slice of string arguments required for command execution
// err: validation or processing error, nil if command is valid
type command struct {
	name string
	data []string
	err  error
}

// represents a single medical visit record
// stores essential information about doctor appointments
// doctor: name of the healthcare provider
// date: timestamp of the appointment in YYYY-MM-DD format
type visit struct {
	doctor string
	date   time.Time
}
