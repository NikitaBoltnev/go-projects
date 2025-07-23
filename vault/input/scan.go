// Package input provides functions for reading user input from the command line.
package input

import (
	"fmt"
)

func PromptData(promt ...string) string {
	for i, line := range promt {
		if i == len(promt)-1 {
			fmt.Printf("%v: ", line)
		} else {
			fmt.Println(line)
		}
	}
	var res string
	fmt.Scanln(&res)
	return res
}
