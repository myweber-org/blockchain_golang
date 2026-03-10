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

func encryptFile(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
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

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	block, err := aes.NewCipher(key)
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

func main() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		fmt.Printf("Generate key failed: %v\n", err)
		return
	}

	inputFile := "plain.txt"
	encryptedFile := "encrypted.bin"
	decryptedFile := "decrypted.txt"

	if err := os.WriteFile(inputFile, []byte("Secret data for encryption test"), 0644); err != nil {
		fmt.Printf("Create test file failed: %v\n", err)
		return
	}

	if err := encryptFile(inputFile, encryptedFile, key); err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	if err := decryptFile(encryptedFile, decryptedFile, key); err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")

	os.Remove(inputFile)
	os.Remove(encryptedFile)
	os.Remove(decryptedFile)
}package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

type Encryptor struct {
	key []byte
}

func NewEncryptor(key string) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}
	return &Encryptor{key: []byte(key)}, nil
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
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

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func main() {
	key := "thisis32bitlongpassphraseimusing!"
	encryptor, err := NewEncryptor(key)
	if err != nil {
		fmt.Printf("Error creating encryptor: %v\n", err)
		os.Exit(1)
	}

	secret := "confidential data"
	encrypted, err := encryptor.Encrypt(secret)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Encrypted: %s\n", encrypted)

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Decrypted: %s\n", decrypted)
}