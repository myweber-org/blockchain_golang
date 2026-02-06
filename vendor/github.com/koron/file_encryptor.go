package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

type Encryptor struct {
	key []byte
}

func NewEncryptor(key string) (*Encryptor, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(decodedKey) != 32 {
		return nil, errors.New("key must be 32 bytes (64 hex characters)")
	}
	return &Encryptor{key: decodedKey}, nil
}

func (e *Encryptor) EncryptFile(inputPath, outputPath string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(e.key)
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

func (e *Encryptor) DecryptFile(inputPath, outputPath string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(e.key)
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

func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Generate key: file_encryptor -genkey")
		fmt.Println("  Encrypt: file_encryptor -encrypt <key> <input> <output>")
		fmt.Println("  Decrypt: file_encryptor -decrypt <key> <input> <output>")
		return
	}

	switch os.Args[1] {
	case "-genkey":
		key, err := generateRandomKey()
		if err != nil {
			fmt.Printf("Error generating key: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated key: %s\n", key)

	case "-encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor -encrypt <key> <input> <output>")
			os.Exit(1)
		}
		encryptor, err := NewEncryptor(os.Args[2])
		if err != nil {
			fmt.Printf("Error creating encryptor: %v\n", err)
			os.Exit(1)
		}
		err = encryptor.EncryptFile(os.Args[3], os.Args[4])
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")

	case "-decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryptor -decrypt <key> <input> <output>")
			os.Exit(1)
		}
		encryptor, err := NewEncryptor(os.Args[2])
		if err != nil {
			fmt.Printf("Error creating encryptor: %v\n", err)
			os.Exit(1)
		}
		err = encryptor.DecryptFile(os.Args[3], os.Args[4])
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully")

	default:
		fmt.Println("Invalid command")
		os.Exit(1)
	}
}