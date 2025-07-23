// Package main implements a simple guestbook web application.
// It allows users to view and add signatures, which are stored in a text file.
// The server runs on localhost:8080.
package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Guestbook holds data to be displayed in the HTML template.
type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

// check logs a fatal error and exits if err is not nil.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// viewHandler serves the guestbook page with all signatures.
func viewHandler(write http.ResponseWriter, request *http.Request) {
	signatures := getStrings("data/signatures.txt")
	html, err := template.ParseFiles("templates/view.html")
	check(err)
	guestbook := Guestbook{
		SignatureCount: len(signatures),
		Signatures:     signatures,
	}
	err = html.Execute(write, guestbook)
	check(err)
}

// formHandler serves the HTML form for adding a new signature.
func formHandler(write http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("templates/form.html")
	check(err)
	err = html.Execute(write, nil)
	check(err)
}

// createHandler saves a new signature from the form and redirects to /guestbook.
func createHandler(write http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("data/signatures.txt", options, os.FileMode(0600))
	check(err)
	_, err = fmt.Fprintln(file, signature)
	check(err)
	err = file.Close()
	check(err)
	http.Redirect(write, request, "/guestbook", http.StatusFound)
}

// getStrings reads all lines from a file and returns them as a slice of strings.
func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}

// main starts an HTTP server and registers three routes.
// The server listens on localhost:8080.
func main() {
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/form", formHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
