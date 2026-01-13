
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

func encryptData(plaintext []byte, passphrase string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encrypted string, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase)
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <passphrase>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputFile := os.Args[2]
	passphrase := os.Args[3]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch operation {
	case "encrypt":
		encrypted, err := encryptData(data, passphrase)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
		outputFile := inputFile + ".enc"
		if err := os.WriteFile(outputFile, []byte(encrypted), 0644); err != nil {
			fmt.Printf("Error writing encrypted file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted data written to %s\n", outputFile)

	case "decrypt":
		decrypted, err := decryptData(string(data), passphrase)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
		outputFile := inputFile + ".dec"
		if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
			fmt.Printf("Error writing decrypted file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted data written to %s\n", outputFile)

	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}