package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const saltSize = 32

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(password, salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	outputData := append(salt, ciphertext...)

	return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, password string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if len(data) < saltSize {
		return errors.New("file too short")
	}

	salt := data[:saltSize]
	ciphertext := data[saltSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output>\n", filepath.Base(os.Args[0]))
		fmt.Println("Password will be read from ENCRYPTION_PASSWORD environment variable")
		os.Exit(1)
	}

	password := os.Getenv("ENCRYPTION_PASSWORD")
	if password == "" {
		fmt.Println("Error: ENCRYPTION_PASSWORD environment variable not set")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		fmt.Printf("Unknown operation: %s\n", operation)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation %s completed successfully\n", operation)
}