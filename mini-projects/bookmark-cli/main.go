// Package main implements a simple CLI bookmark manager.
// Users can view, add, and delete bookmarks (key-value pairs).
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// bookmarkMap is a type alias for map[string]string to store bookmarks.
type bookmarkMap = map[string]string

// main runs the CLI bookmark manager.
// It displays a menu and processes user input in a loop until exit.
func main() {
	bookmarks := make(bookmarkMap)
Menu:
	for {
		printMenu()
		input := getUserInput()

		switch input {
		case 1:
			outputBookmarks(bookmarks)
		case 2:
			inputBookmark(bookmarks)
		case 3:
			deleteBookmark(bookmarks)
		case 4:
			fmt.Println("Goodbye!")
			break Menu
		default:
			fmt.Println("Invalid input. Please enter a number between 1 and 4.")

		}
	}

}

// printMenu displays the main menu options.
func printMenu() {
	fmt.Println("\nChoose an action:")
	fmt.Println("1 - View bookmarks")
	fmt.Println("2 - Add a bookmark")
	fmt.Println("3 - Delete a bookmark")
	fmt.Println("4 - Exit")

}

// getUserInput reads user input.
func getUserInput() int {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0
	}
	input = strings.TrimSpace(input)
	number, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0
	}

	return int(number)
}

// outputBookmarks prints all stored bookmarks.
func outputBookmarks(bookmarks bookmarkMap) {
	if len(bookmarks) == 0 {
		fmt.Println("No bookmarks yet.")
	}
	for key, value := range bookmarks {
		fmt.Println(key, ":", value)
	}
}

// inputBookmark prompts the user to enter a key and value, then adds it.
func inputBookmark(bookmarks bookmarkMap) {
	fmt.Print("Enter key: ")
	key := getUserString()
	fmt.Print("Enter value: ")
	value := getUserString()
	bookmarks[key] = value
	fmt.Printf("Bookmark '%s' added.\n", key)
}

// getUserString reads a single line of input and returns it as a trimmed string.
func getUserString() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSpace(input)
	return input
}

// deleteBookmark prompts for a key and deletes the corresponding bookmark.
func deleteBookmark(bookmarks bookmarkMap) {
	fmt.Print("Enter key to delete: ")
	key := getUserString()

	if _, exists := bookmarks[key]; exists {
		delete(bookmarks, key)
		fmt.Printf("Bookmark '%s' deleted.\n", key)
	} else {
		fmt.Printf("Bookmark '%s' not found.\n", key)
	}
}
