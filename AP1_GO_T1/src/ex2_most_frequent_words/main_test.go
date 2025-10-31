package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal",
			input:    "aa bb cc aa cc cc cc aa ab ac bb\n3",
			expected: "cc aa bb\n",
		},
		{
			name:     "normal with lexicographically",
			input:    "aa bb cc aa cc cc cc aa ab ac bb ff ff\n4",
			expected: "cc aa bb ff\n",
		},
		{
			name:     "empty list",
			input:    "\n5",
			expected: "",
		},
		{
			name:     "list contains fewer unique words than K",
			input:    "aa bb cc aa cc cc cc aa ab ac bb\n2",
			expected: "cc aa\n",
		},
		{
			name:     "list contains fewer unique words than K with lexicographically",
			input:    "aa bb cc aa cc cc cc aa ab ac bb aa\n2",
			expected: "aa cc\n",
		},
		{
			name:     "new",
			input:    "aa bb cc\n20",
			expected: "aa bb cc\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			runProgram(strings.NewReader(tt.input), &buf)

			got := buf.String()

			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}

		})
	}
}
