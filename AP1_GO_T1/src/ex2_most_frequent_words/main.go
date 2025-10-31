package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	//running a program with standard input/output streams
	runProgram(os.Stdin, os.Stdout)
}

func runProgram(input io.Reader, output io.Writer) {
	wordFreqs, num := parseInput(input)

	uniqueWords := getUniqueWords(wordFreqs)
	frequencies := getWordFrequencies(uniqueWords, wordFreqs)
	indexes := sorting(uniqueWords, frequencies)

	uniqueCount := len(uniqueWords)
	if uniqueCount < num {
		num = uniqueCount
	}

	indexes = indexes[:num]
	res := make([]string, 0, num)
	for _, v := range indexes {
		res = append(res, uniqueWords[v])
	}

	outputResult(res, num, output)
}

func outputResult(res []string, num int, output io.Writer) {
	if num == 0 {
		return
	}

	for i := range num {
		fmt.Fprint(output, res[i])

		if i < num-1 {
			fmt.Fprint(output, " ")
		}
	}
	fmt.Fprint(output, "\n")
}

// returns slice of indices that define the output order of words
func sorting(uniqueWords []string, frequencies []int) []int {
	indexes := make([]int, len(uniqueWords))
	for i := range indexes {
		indexes[i] = i
	}

	sort.Slice(indexes, func(i, j int) bool {
		// if two words have the same frequency
		if frequencies[indexes[i]] == frequencies[indexes[j]] {
			return uniqueWords[indexes[i]] < uniqueWords[indexes[j]]
		}
		// otherwise sort by frequency in descending order
		return frequencies[indexes[i]] > frequencies[indexes[j]]
	})

	return indexes
}

func getWordFrequencies(uniqueWords []string, wordFreqs map[string]int) []int {
	frequencies := make([]int, len(uniqueWords))

	for i, word := range uniqueWords {
		frequencies[i] = wordFreqs[word]
	}

	return frequencies
}

func getUniqueWords(wordFreqs map[string]int) []string {
	uniqueWords := make([]string, 0, len(wordFreqs))

	for n := range wordFreqs {
		uniqueWords = append(uniqueWords, n)
	}
	sort.Strings(uniqueWords)

	return uniqueWords
}

func parseInput(input io.Reader) (map[string]int, int) {
	var frequencies int
	var err error
	wordFreqs := make(map[string]int)

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		temp := scanner.Text()

		frequencies, err = strconv.Atoi(temp)
		if err == nil {
			break
		}
		wordFreqs[temp]++
	}

	return wordFreqs, frequencies
}
