
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
)

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := deriveKey(passphrase, salt)

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
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	outputData := append(salt, ciphertext...)

	return os.WriteFile(outputPath, outputData, 0644)
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	ciphertextWithSalt, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if len(ciphertextWithSalt) < 16 {
		return errors.New("file too short")
	}

	salt := ciphertextWithSalt[:16]
	ciphertext := ciphertextWithSalt[16:]

	key := deriveKey(passphrase, salt)

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
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	action := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]
	passphrase := os.Args[4]

	var err error
	switch action {
	case "encrypt":
		err = encryptFile(inputFile, outputFile, passphrase)
	case "decrypt":
		err = decryptFile(inputFile, outputFile, passphrase)
	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation completed successfully: %s -> %s\n", inputFile, outputFile)
}