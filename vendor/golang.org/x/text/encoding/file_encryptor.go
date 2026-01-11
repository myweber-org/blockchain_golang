
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

type FileEncryptor struct {
	key []byte
}

func NewFileEncryptor(key string) (*FileEncryptor, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid hex key: %v", err)
	}
	if len(decodedKey) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes for AES-256")
	}
	return &FileEncryptor{key: decodedKey}, nil
}

func (fe *FileEncryptor) EncryptFile(inputPath, outputPath string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %v", err)
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation failed: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	return nil
}

func (fe *FileEncryptor) DecryptFile(inputPath, outputPath string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file failed: %v", err)
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return fmt.Errorf("cipher creation failed: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM creation failed: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: %v", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	return nil
}

func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("key generation failed: %v", err)
	}
	return hex.EncodeToString(key), nil
}

func main() {
	key, err := generateRandomKey()
	if err != nil {
		fmt.Printf("Key generation error: %v\n", err)
		return
	}
	fmt.Printf("Generated key: %s\n", key)

	encryptor, err := NewFileEncryptor(key)
	if err != nil {
		fmt.Printf("Encryptor creation error: %v\n", err)
		return
	}

	testData := []byte("Confidential data requiring encryption")
	testFile := "test_data.txt"
	encryptedFile := "encrypted.dat"
	decryptedFile := "decrypted.txt"

	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		fmt.Printf("Test file creation error: %v\n", err)
		return
	}
	defer os.Remove(testFile)
	defer os.Remove(encryptedFile)
	defer os.Remove(decryptedFile)

	if err := encryptor.EncryptFile(testFile, encryptedFile); err != nil {
		fmt.Printf("Encryption error: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	if err := encryptor.DecryptFile(encryptedFile, decryptedFile); err != nil {
		fmt.Printf("Decryption error: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		fmt.Printf("Read decrypted file error: %v\n", err)
		return
	}

	if string(decryptedData) == string(testData) {
		fmt.Println("Encryption/decryption verification successful")
	} else {
		fmt.Println("Encryption/decryption verification failed")
	}
}