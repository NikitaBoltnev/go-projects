// Package files provides a simple JSON file storage implementation.
// It supports reading from and writing to a file with error handling and success feedback.
package files

import (
	"os"
	"vault/output"

	"github.com/fatih/color"
)

type JsonDb struct {
	filename string
}

func NewJsonDb(name string) *JsonDb {
	return &JsonDb{
		filename: name,
	}
}

// Write saves the given content to the file.
func (db *JsonDb) Write(content []byte) {
	file, err := os.Create(db.filename)
	if err != nil {
		output.PrintError(err)
		return
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		output.PrintError(err)
		return
	}
	color.Green("Write successful.")
}

// Read reads and returns the content of the file as a byte slice.
func (db *JsonDb) Read() ([]byte, error) {
	data, err := os.ReadFile(db.filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}
