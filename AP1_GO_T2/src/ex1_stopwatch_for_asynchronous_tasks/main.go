package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type stat struct {
	ID       int
	workTime int
}

func main() {
	// N: number of goroutines to launch
	// M: maximum sleep time in milliseconds for each goroutine
	N := flag.Int("N", 1, "сount goroutines")
	M := flag.Int("M", 10, "goroutine sleep time")
	flag.Parse()

	if err := check(*N); err != nil {
		fmt.Println("N must be greater than 0")
		return
	}
	if err := check(*M); err != nil {
		fmt.Println("M must be greater than 0")
		return
	}

	var wg sync.WaitGroup
	res := make([]stat, 0, *N)

	for i := range *N {
		wg.Add(1)
		randomTime := rand.Intn(*M) + 1 // generate random sleep time
		res = append(res, stat{i, randomTime})

		go func(randomTime int) {
			defer wg.Done()
			time.Sleep(time.Duration(randomTime) * time.Millisecond)
		}(randomTime)
	}

	wg.Wait()

	sort.Slice(res, func(i, j int) bool {
		return res[i].workTime > res[j].workTime
	})

	output(res)
}

func check(n int) error {
	if n <= 0 {
		return errors.New("invalid number")
	}
	return nil
}

func output(res []stat) {
	for _, v := range res {
		fmt.Print("ID: ", v.ID, " time: ", v.workTime, "\n")
	}
}
