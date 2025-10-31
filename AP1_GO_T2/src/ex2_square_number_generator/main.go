package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
)

func main() {
	K := flag.Int("K", 0, "start of number range")
	N := flag.Int("N", 0, "end of number range")

	flag.Parse()

	err := check(*K, *N)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch1 := generator(*K, *N)
	ch2 := squaring(ch1)

	for res := range ch2 {
		fmt.Println(res)
	}

}

// creates a channel and produces numbers from k to n (inclusive)
// returns a read-only channel to prevent accidental writing by consumers
func generator(k, n int) <-chan int {
	ch1 := make(chan int)

	go func(ch1 chan int, k, n int) {
		defer close(ch1)
		for i := k; i <= n; i++ {
			ch1 <- i
		}
	}(ch1, k, n)

	return ch1
}

// reads numbers from input channel, squares them, and sends to output channel
// returns read-only output channel for safe consumption
func squaring(ch1 <-chan int) <-chan int {
	ch2 := make(chan int)

	go func(ch1 <-chan int, ch2 chan int) {
		defer close(ch2)
		for num := range ch1 {
			ch2 <- int(math.Pow(float64(num), 2))
		}
	}(ch1, ch2)

	return ch2
}

// validates that the number range is valid
func check(k, n int) error {
	if n < k {
		return errors.New("N must be greater than or equal to K")
	}

	return nil
}
