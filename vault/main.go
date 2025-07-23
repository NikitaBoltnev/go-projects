// Package main implements a simple password manager CLI.
// It allows users to create, search, and delete account entries stored in an encrypted vault.
package main

import (
	"fmt"
	"strings"
	"vault/account"
	"vault/encrypter"
	"vault/files"
	"vault/input"
	"vault/output"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var menu = map[string]func(*account.VaultWithDb){
	"1": createAccount,
	"2": findAccountByUrl,
	"3": findAccountByLogin,
	"4": deleteAccount,
}

var menuVariants = []string{
	"1. Create account",
	"2. Find account by URL",
	"3. Find account by login",
	"4. Delete account",
	"5. Exit",
	"Select an option",
}

// main is the entry point of the password manager.
// It loads environment variables, initializes the encrypted vault,
// and runs the main menu loop for user interaction.
func main() {
	err := godotenv.Load()
	if err != nil {
		output.PrintError("Environment file (.env) not found")
	}
	fmt.Println("__Password Manager__")
	vault := account.NewVault(files.NewJsonDb("data.vault"), *encrypter.NewEncrypter())

Menu:
	for {
		variant := input.PromptData(menuVariants...)
		menuFunc := menu[variant]
		if menuFunc == nil {
			break Menu
		}
		menuFunc(vault)
	}
}

// createAccount prompts the user for login, password, and URL,
// creates a new account, and adds it to the vault.
func createAccount(vault *account.VaultWithDb) {
	login := input.PromptData("Enter login")
	password := input.PromptData("Enter password or leave empty to generate one")
	url := input.PromptData("Enter URL")

	myAccount, err := account.NewAccount(login, password, url)
	if err != nil {
		output.PrintError("Invalid URL or login format.")
		return
	}

	vault.AddAccount(*myAccount)
}

// findAccountByUrl prompts the user for a URL and searches for matching accounts.
func findAccountByUrl(vault *account.VaultWithDb) {
	url := input.PromptData("Enter URL to search")
	accounts := vault.FindAccounts(url, func(acc account.Account, str string) bool {
		return strings.Contains(acc.Url, str)
	})
	outputResult(&accounts)
}

// findAccountByLogin prompts the user for a login and searches for matching accounts.
func findAccountByLogin(vault *account.VaultWithDb) {
	login := input.PromptData("Enter login to search")
	accounts := vault.FindAccounts(login, func(acc account.Account, str string) bool {
		return strings.Contains(acc.Login, str)
	})
	outputResult(&accounts)
}

// outputResult prints the list of found accounts or an error if none are found.
func outputResult(accounts *[]account.Account) {
	if len(*accounts) == 0 {
		output.PrintError("No accounts found.")
	}
	for _, account := range *accounts {
		account.Output()
	}
}

// deleteAccount prompts the user for a URL and deletes all matching accounts from the vault.
func deleteAccount(vault *account.VaultWithDb) {
	url := input.PromptData("Enter URL to delete")
	isDeleted := vault.DeleteAccountsByUrl(url)
	if isDeleted {
		color.Green("Deleted")
	} else {
		output.PrintError("Not found")
	}
}
