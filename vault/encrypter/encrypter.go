// Package encrypter provides AES-GCM encryption and decryption utilities.
package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
)

type Encrypter struct {
	Key string
}

// NewEncrypter creates a new Encrypter instance using the KEY from environment variables.
func NewEncrypter() *Encrypter {
	key := os.Getenv("KEY")
	if key == "" {
		panic("KEY environment variable not set")
	}
	return &Encrypter{
		Key: key,
	}
}

// Encrypt encrypts the given plaintext byte slice using AES-GCM.
func (enc *Encrypter) Encrypt(plainStr []byte) []byte {
	keyBytes, err := hex.DecodeString(enc.Key)
	if err != nil {
		panic(err.Error())
	}
	block, err := aes.NewCipher([]byte(keyBytes))
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		panic(err.Error())
	}
	return aesGCM.Seal(nonce, nonce, plainStr, nil)
}

// Decrypt decrypts the given ciphertext (with prepended nonce) using AES-GCM.
func (enc *Encrypter) Decrypt(encryptedStr []byte) []byte {
	keyBytes, err := hex.DecodeString(enc.Key)
	if err != nil {
		panic(err.Error())
	}
	block, err := aes.NewCipher([]byte(keyBytes))
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonseSize := aesGCM.NonceSize()
	nonce, cipherText := encryptedStr[:nonseSize], encryptedStr[nonseSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		panic(err.Error())
	}
	return plainText
}
