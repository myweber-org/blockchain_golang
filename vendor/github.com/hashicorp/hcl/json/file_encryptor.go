
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

type Encryptor struct {
	key []byte
}

func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}
	return &Encryptor{key: key}, nil
}

func (e *Encryptor) EncryptFile(inputPath, outputPath string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func (e *Encryptor) DecryptFile(inputPath, outputPath string) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}
	return key, nil
}package main

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
	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	original := []byte("Sensitive data requiring protection")
	fmt.Printf("Original: %s\n", original)

	encrypted, err := encryptData(original, key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := decryptData(encrypted, key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}