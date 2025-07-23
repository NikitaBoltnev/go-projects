// Package account provides functionality for managing encrypted account storage.
package account

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
	"vault/encrypter"
	"vault/output"

	"github.com/fatih/color"
)

type Db interface {
	Read() ([]byte, error)
	Write([]byte)
}

type Vault struct {
	Accounts  []Account `json:"accounts"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VaultWithDb struct {
	Vault
	db  Db
	enc encrypter.Encrypter
}

// NewVault initializes a vault from an encrypted file.
func NewVault(db Db, enc encrypter.Encrypter) *VaultWithDb {
	file, err := db.Read()
	if err != nil {
		return &VaultWithDb{
			Vault: Vault{
				Accounts:  []Account{},
				UpdatedAt: time.Now(),
			},
			db:  db,
			enc: enc,
		}
	}
	data := enc.Decrypt(file)
	var vault Vault
	err = json.Unmarshal(data, &vault)
	if err != nil {
		output.PrintError("Failed to parse data.vault file")
		return &VaultWithDb{
			Vault: Vault{
				Accounts:  []Account{},
				UpdatedAt: time.Now(),
			},
			db:  db,
			enc: enc,
		}
	}
	color.Cyan("Found %d accounts", len(vault.Accounts))
	return &VaultWithDb{
		Vault: vault,
		db:    db,
		enc:   enc,
	}
}

// AddAccount adds a new account to the vault and saves it.
func (vault *VaultWithDb) AddAccount(acc Account) {
	vault.Accounts = append(vault.Accounts, acc)
	vault.save()
}

// FindAccounts searches for accounts using a custom matching function.
func (vault *VaultWithDb) FindAccounts(str string, checker func(Account, string) bool) []Account {
	var accounts []Account
	for _, account := range vault.Accounts {
		isMatched := checker(account, str)
		if isMatched {
			accounts = append(accounts, account)
		}
	}
	return accounts
}

// DeleteAccountsByUrl removes all accounts whose URL contains the given string.
func (vault *VaultWithDb) DeleteAccountsByUrl(url string) bool {
	var accounts []Account
	isDeleted := false
	for _, account := range vault.Accounts {
		isMatched := strings.Contains(account.Url, url)
		if !isMatched {
			accounts = append(accounts, account)
			continue
		}
		isDeleted = true
	}

	if isDeleted {
		vault.Accounts = accounts
		vault.save()
	}

	return isDeleted
}

// ToBytes serializes the vault to JSON.
func (vault *Vault) ToBytes() ([]byte, error) {
	file, err := json.Marshal(vault)
	if err != nil {
		return nil, errors.New("INVALID_JSON")
	}
	return file, nil
}

// save encrypts and writes the current vault state to the database.
func (vault *VaultWithDb) save() {
	vault.UpdatedAt = time.Now()
	data, err := vault.Vault.ToBytes()
	if err != nil {
		output.PrintError("Failed to serialize vault")
	}
	encData := vault.enc.Encrypt(data)
	vault.db.Write(encData)
}
