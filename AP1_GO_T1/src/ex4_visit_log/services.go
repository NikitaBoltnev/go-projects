package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

func distributor(wg *sync.WaitGroup, inputCh chan string, commandCh chan command) {
	defer wg.Done()
	defer close(commandCh)

	for inputLine := range inputCh {

		var err error

		var command command

		inputLine = strings.TrimSpace(inputLine)
		inputLine = strings.ToLower(inputLine)
		switch {
		case inputLine == "save":
			command.name = "save"
			handleGetData(&command, inputCh, 3) // collect patient, doctor, and date
			err = checkData(&command)

		case inputLine == "gethistory":
			command.name = "getHistory"
			handleGetData(&command, inputCh, 1) // collect patient name only
			err = checkData(&command)

		case inputLine == "getlastvisit":
			command.name = "getLastVisit"
			handleGetData(&command, inputCh, 2) // collect patient and doctor
			err = checkData(&command)

		default:
			err = errors.New("invalid command")
		}

		if err != nil {
			command.err = err
		}

		commandCh <- command
	}

}

func checkData(command *command) error {
	var err error
	if command.name == "save" {
		// validate all three fields for save command
		patientName := strings.Fields(command.data[0])
		doctor := command.data[1]
		date := command.data[2]

		if err = checkDataPatientName(patientName); err != nil {
			return err
		}
		if err = checkDataDoctor(doctor); err != nil {
			return err
		}
		if err = checkDataDate(date); err != nil {
			return err
		}
	}

	if command.name == "getHistory" {
		// validate patient name for history query
		patientName := strings.Fields(command.data[0])
		if err = checkDataPatientName(patientName); err != nil {
			return err
		}
	}

	if command.name == "getLastVisit" {
		// validate both patient and doctor for last visit query
		patientName := strings.Fields(command.data[0])
		doctor := command.data[1]
		if err = checkDataPatientName(patientName); err != nil {
			return err
		}
		if err = checkDataDoctor(doctor); err != nil {
			return err
		}
	}

	return nil
}

// validates doctor name format
// ensures name is not empty and contains only letters and hyphens
func checkDataDoctor(doctor string) error {
	var err error
	if doctor == "" {
		err = errors.New("empty doctor name")
		return err
	}
	for _, v := range doctor {
		if !unicode.IsLetter(v) && v != '-' {
			err = errors.New("wrong doctor name")
		}
	}

	return err
}

// validates date string format
// requires exact YYYY-MM-DD format using reference date 2006-01-02
func checkDataDate(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		err = errors.New("invalid date format")
	}
	return err
}

// validates patient name structure
// requires either 2 (first and last) or 3 (first, middle, last) name components
func checkDataPatientName(patientName []string) error {
	var err error

	if len(patientName) != 3 && len(patientName) != 2 {
		err = errors.New("invalid name entered")
	}

	return err
}

func handleGetData(command *command, inputCh chan string, step int) {
	for range step {
		data := <-inputCh
		data = strings.TrimSpace(data)
		command.data = append(command.data, data)
	}
}

func input(wg *sync.WaitGroup, inputCh chan string) {
	defer wg.Done()
	defer close(inputCh)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		inputLine := scanner.Text()

		if inputLine == "exit" {
			break
		}

		inputCh <- inputLine

	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
	}

}
