package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	errInvalidInput   = errors.New("Invalid input")
	errDivisionByZero = errors.New("Division by zero")
)

// available arithmetic operations
var operation = map[string]func(float64, float64) (float64, error){
	"+": func(a, b float64) (float64, error) { return a + b, nil },
	"-": func(a, b float64) (float64, error) { return a - b, nil },
	"*": func(a, b float64) (float64, error) { return a * b, nil },
	"/": func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, errDivisionByZero
		}
		return a / b, nil
	},
}

func main() {
	var validCalculation bool // track if calculation succeeded
	var res float64
	var err error

	leftOperand := inputOperand("left operand")
	op := inputOperation()

	for !validCalculation {
		rightOperand := inputOperand("right operand")
		operationFunc := operation[op]
		res, err = operationFunc(leftOperand, rightOperand)

		// if division by zero occurred, request right operand again
		if err != nil {
			fmt.Println(errInvalidInput)
			fmt.Println(err)
			continue
		}
		validCalculation = true
	}

	if op == "/" { //format division result
		fmt.Printf("Result: %.3f\n", res)
	} else {
		fmt.Printf("Result: %g\n", res)
	}
}

func inputOperation() string {
	var success bool
	var operation string
	var err error

	// enter the operation until it is correct
	for !success {
		fmt.Println("Input operation:")
		operation, err = reader()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Read error:", err)
			continue
		}

		if !strings.Contains("+-*/", operation) {
			fmt.Println(errInvalidInput)
			fmt.Println("Use: +, -, *, or /")
		} else {
			success = true
		}
	}

	return operation
}

func inputOperand(str string) float64 {
	var num float64
	var success bool

	// enter the operand until it is correct
	for !success {
		fmt.Println("Input " + str + ":")
		operand, err := reader()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Read error:", err)
			continue
		}

		num, err = strconv.ParseFloat(operand, 64)
		if err != nil {
			fmt.Println(errInvalidInput)
			fmt.Println("Examples: 198, 12.345, -0.987")
		} else {
			success = true
		}
	}

	return num
}

func reader() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	return input, nil
}
