package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

func encryptData(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_encryption_tool.go <command>")
		fmt.Println("Commands: generate-key, encrypt <text>, decrypt <ciphertext>")
		return
	}

	command := os.Args[1]

	switch command {
	case "generate-key":
		key, err := generateKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			return
		}
		fmt.Printf("Generated key: %x\n", key)

	case "encrypt":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run file_encryption_tool.go encrypt <text> <key_hex>")
			return
		}
		text := []byte(os.Args[2])
		var key []byte
		fmt.Sscanf(os.Args[3], "%x", &key)

		ciphertext, err := encryptData(text, key)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			return
		}
		fmt.Printf("Ciphertext: %x\n", ciphertext)

	case "decrypt":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run file_encryption_tool.go decrypt <ciphertext_hex> <key_hex>")
			return
		}
		var ciphertext []byte
		fmt.Sscanf(os.Args[2], "%x", &ciphertext)
		var key []byte
		fmt.Sscanf(os.Args[3], "%x", &key)

		plaintext, err := decryptData(ciphertext, key)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			return
		}
		fmt.Printf("Plaintext: %s\n", plaintext)

	default:
		fmt.Println("Unknown command")
	}
}