package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	visits := make(map[string][]visit) // map to store patient visits data
	var wg sync.WaitGroup
	// channel for passing user input from input goroutine
	ch := make(chan string)
	// channel for passing commands to be executed from distributor goroutine
	reqCh := make(chan command)

	wg.Add(2)
	go input(&wg, ch)
	go distributor(&wg, ch, reqCh)

	// process commands from the channel until it's closed
	for cmd := range reqCh {
		if cmd.err != nil {
			fmt.Println(cmd.err.Error())
		} else {
			executeCommand(visits, cmd)
		}
	}
	wg.Wait()
}

func executeCommand(visits map[string][]visit, cmd command) {
	var res any
	var err error
	switch cmd.name {
	case "save":
		setVisit(visits, cmd)
		return

	case "getHistory":
		res, err = getHistory(visits, cmd)

	case "getLastVisit":
		res, err = getLastVisit(visits, cmd)

	default:
		err = errors.New("unknown command")
	}

	if err != nil {
		fmt.Println(err.Error())
	} else {
		printResult(res)
	}
}

// formats and displays command execution results based on data type
func printResult(result any) {
	switch visitRecord := result.(type) {
	case time.Time:
		fmt.Println(visitRecord.Format("2006-01-02"))

	case []visit:
		for _, visit := range visitRecord {
			fmt.Printf("%s %s\n", visit.doctor, visit.date.Format("2006-01-02"))
		}

	default:
		fmt.Println("unknown data type")
	}
}
