package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func encryptData(plaintext []byte, key []byte) (string, error) {
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

func decryptData(encrypted string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
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
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt>")
		return
	}

	action := os.Args[1]
	key, _ := generateKey()

	switch action {
	case "encrypt":
		plaintext := []byte("Sensitive data to protect")
		encrypted, err := encryptData(plaintext, key)
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			return
		}
		fmt.Printf("Encrypted: %s\n", encrypted)
		fmt.Printf("Key (base64): %s\n", base64.StdEncoding.EncodeToString(key))

	case "decrypt":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run file_encryption_tool.go decrypt <encrypted_data> <key_base64>")
			return
		}
		encryptedData := os.Args[2]
		keyBase64 := os.Args[3]

		key, err := base64.StdEncoding.DecodeString(keyBase64)
		if err != nil {
			fmt.Printf("Invalid key: %v\n", err)
			return
		}

		decrypted, err := decryptData(encryptedData, key)
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			return
		}
		fmt.Printf("Decrypted: %s\n", decrypted)

	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
	}
}