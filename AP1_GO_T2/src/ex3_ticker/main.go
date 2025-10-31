package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	K := flag.Uint("K", 1, "Tick step in seconds")
	flag.Parse()

	err := check(*K)
	if err != nil {
		return
	}

	// channel is used for graceful shutdown signaling
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)

	go stopper(&wg, done)
	go ticker(&wg, done, *K)

	wg.Wait()

}

// generates periodic ticks at the specified interval
func ticker(wg *sync.WaitGroup, done <-chan bool, step uint) {
	defer wg.Done()

	i := uint(1)   // tick counter
	sec := uint(0) // total elapsed seconds

	for {

		select {
		case <-done:
			fmt.Println("Termination")
			return

		default:
			time.Sleep(time.Second * time.Duration(step))
			sec += step
			fmt.Printf("Tick %d, since %d\n", i, sec)
			i++
		}

	}

}

// listens for OS termination signals (SIGTERM, SIGINT)
func stopper(wg *sync.WaitGroup, done chan<- bool) {
	sigCh := make(chan os.Signal, 1)

	defer wg.Done()
	defer close(done)

	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh

	signal.Stop(sigCh)
}

func check(k uint) error {
	if k == 0 {
		return errors.New("K cannot be zero")
	}
	return nil
}
