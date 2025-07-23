// Package output provides utilities for printing formatted error messages to the console.
package output

import (
	"github.com/fatih/color"
)

func PrintError(value any) {
	switch t := value.(type) {
	case string:
		color.Red(t)
	case int:
		color.Red("Error code: %d", t)
	case error:
		color.Red(t.Error())
	default:
		color.Red("Unknown error")
	}
}
