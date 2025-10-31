package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// disable log prefixes for cleaner output
	log.SetFlags(0)
	arr1, err := readNumbers()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid input")
		os.Exit(0)
	}
	arr2, err := readNumbers()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid input")
		os.Exit(0)
	}

	res := findIntersection(arr1, arr2)
	if len(res) > 0 {
		printNumbers(res)
	} else {
		fmt.Println("Empty intersection")
	}
}

func printNumbers(res []int) {
	for n, v := range res {
		fmt.Print(v)

		if n != len(res)-1 {
			fmt.Print(" ")
		}

	}
	fmt.Println()
}

// returns common elements between two arrays
// preserves order from first array and removes duplicates
// uses map for O(1) lookups from second array
func findIntersection(arr1, arr2 []int) []int {
	set := make(map[int]bool)
	for _, num := range arr2 {
		set[num] = true
	}

	var res []int

	for _, num := range arr1 {
		if set[num] {
			res = append(res, num)
			set[num] = false
		}
	}

	return res
}

// reads a line from stdin and parses space-separated integers
// returns error if any token cannot be converted to integer
func readNumbers() ([]int, error) {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	var numbers []int

	for _, v := range strings.Fields(input) {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, num)
	}

	return numbers, nil
}
