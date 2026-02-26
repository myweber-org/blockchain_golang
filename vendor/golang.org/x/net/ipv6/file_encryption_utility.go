package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize      = 16
	nonceSize     = 12
	keyIterations = 100000
	keyLength     = 32
)

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, keyIterations, keyLength, sha256.New)
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

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	outputData := append(salt, nonce...)
	outputData = append(outputData, ciphertext...)

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

	if len(ciphertext) < saltSize+nonceSize {
		return errors.New("file too short to contain valid encrypted data")
	}

	salt := ciphertext[:saltSize]
	nonce := ciphertext[saltSize : saltSize+nonceSize]
	actualCiphertext := ciphertext[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func generateRandomFile(path string, size int) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Errorf("failed to generate random data: %w", err)
	}

	if _, err := file.Write(buf); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_encryption_utility <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  encrypt <input> <output> <password>")
		fmt.Println("  decrypt <input> <output> <password>")
		fmt.Println("  test <size> <password>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "encrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryption_utility encrypt <input> <output> <password>")
			os.Exit(1)
		}
		if err := encryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Encryption completed successfully")

	case "decrypt":
		if len(os.Args) != 5 {
			fmt.Println("Usage: file_encryption_utility decrypt <input> <output> <password>")
			os.Exit(1)
		}
		if err := decryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Decryption completed successfully")

	case "test":
		if len(os.Args) != 4 {
			fmt.Println("Usage: file_encryption_utility test <size> <password>")
			os.Exit(1)
		}

		var size int
		if _, err := fmt.Sscanf(os.Args[2], "%d", &size); err != nil {
			fmt.Printf("Invalid size: %v\n", err)
			os.Exit(1)
		}

		password := os.Args[3]
		tempInput := "test_input.bin"
		tempEncrypted := "test_encrypted.bin"
		tempDecrypted := "test_decrypted.bin"

		defer func() {
			os.Remove(tempInput)
			os.Remove(tempEncrypted)
			os.Remove(tempDecrypted)
		}()

		fmt.Printf("Generating random file of size %d bytes...\n", size)
		if err := generateRandomFile(tempInput, size); err != nil {
			fmt.Printf("Failed to generate test file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Encrypting file...")
		if err := encryptFile(tempInput, tempEncrypted, password); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Decrypting file...")
		if err := decryptFile(tempEncrypted, tempDecrypted, password); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}

		original, _ := os.ReadFile(tempInput)
		decrypted, _ := os.ReadFile(tempDecrypted)

		if hex.EncodeToString(original) == hex.EncodeToString(decrypted) {
			fmt.Println("Test passed: Original and decrypted files match")
		} else {
			fmt.Println("Test failed: Original and decrypted files do not match")
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}