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

const (
	saltSize   = 16
	nonceSize  = 12
	keyIter    = 100000
)

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	for i := 0; i < keyIter; i++ {
		combined := append(hash[:], salt...)
		hash = sha256.Sum256(combined)
	}
	return hash[:]
}

func encryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("salt generation failed: %w", err)
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode failed: %w", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %w", err)
	}
	defer inputFile.Close()

	plaintext, err := io.ReadAll(inputFile)
	if err != nil {
		return fmt.Errorf("file read failed: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %w", err)
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(salt); err != nil {
		return fmt.Errorf("salt write failed: %w", err)
	}
	if _, err := outputFile.Write(nonce); err != nil {
		return fmt.Errorf("nonce write failed: %w", err)
	}
	if _, err := outputFile.Write(ciphertext); err != nil {
		return fmt.Errorf("ciphertext write failed: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("input file open failed: %w", err)
	}
	defer inputFile.Close()

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(inputFile, salt); err != nil {
		return fmt.Errorf("salt read failed: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(inputFile, nonce); err != nil {
		return fmt.Errorf("nonce read failed: %w", err)
	}

	ciphertext, err := io.ReadAll(inputFile)
	if err != nil {
		return fmt.Errorf("ciphertext read failed: %w", err)
	}

	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode failed: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed - incorrect passphrase or corrupted file")
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("output file creation failed: %w", err)
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(plaintext); err != nil {
		return fmt.Errorf("plaintext write failed: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]
	passphrase := os.Args[4]

	switch operation {
	case "encrypt":
		if err := encryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")
	case "decrypt":
		if err := decryptFile(inputPath, outputPath, passphrase); err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully")
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}