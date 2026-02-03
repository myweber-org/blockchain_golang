package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func encryptFile(inputPath, outputPath, keyString string) error {
	key, _ := hex.DecodeString(keyString)
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
	return os.WriteFile(outputPath, ciphertext, 0644)
}

func decryptFile(inputPath, outputPath, keyString string) error {
	key, _ := hex.DecodeString(keyString)
	ciphertext, err := os.ReadFile(inputPath)
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

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}

func generateKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt|keygen> <input> <output> [key]")
		return
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	switch operation {
	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Key required for encryption")
			return
		}
		key := os.Args[4]
		if err := encryptFile(inputPath, outputPath, key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
		} else {
			fmt.Println("Encryption completed")
		}
	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Key required for decryption")
			return
		}
		key := os.Args[4]
		if err := decryptFile(inputPath, outputPath, key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
		} else {
			fmt.Println("Decryption completed")
		}
	case "keygen":
		fmt.Printf("Generated key: %s\n", generateKey())
	default:
		fmt.Println("Invalid operation")
	}
}