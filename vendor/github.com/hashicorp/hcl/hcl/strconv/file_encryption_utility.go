
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

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	plaintext, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(salt); err != nil {
		return err
	}
	if _, err := outputFile.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	data, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	if len(data) < 16 {
		return errors.New("file too short")
	}

	salt := data[:16]
	ciphertext := data[16:]

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

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(plaintext); err != nil {
		return err
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
		err := encryptFile(inputPath, outputPath, passphrase)
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File encrypted successfully")
	case "decrypt":
		err := decryptFile(inputPath, outputPath, passphrase)
		if err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully")
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
		os.Exit(1)
	}
}