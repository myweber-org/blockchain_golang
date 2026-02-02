
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

const saltSize = 16

func deriveKey(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, password string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	key := deriveKey(password, salt)

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	outputData := append(salt, ciphertext...)

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, password string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if len(ciphertext) < saltSize {
		return errors.New("file too short to contain salt")
	}

	salt := ciphertext[:saltSize]
	key := deriveKey(password, salt)
	ciphertext = ciphertext[saltSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed - wrong password or corrupted file")
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <password>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	password := os.Args[4]

	var err error
	switch operation {
	case "encrypt":
		err = encryptFile(inputPath, outputPath, password)
	case "decrypt":
		err = decryptFile(inputPath, outputPath, password)
	default:
		fmt.Println("Operation must be 'encrypt' or 'decrypt'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Operation completed successfully")
}