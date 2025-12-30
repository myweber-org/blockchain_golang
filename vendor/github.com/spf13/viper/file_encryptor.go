
package main

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

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

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
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
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

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func saveKeyToFile(key []byte, filename string) error {
	encoded := base64.StdEncoding.EncodeToString(key)
	return os.WriteFile(filename, []byte(encoded), 0600)
}

func loadKeyFromFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func main() {
	key, err := generateKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	err = saveKeyToFile(key, "encryption.key")
	if err != nil {
		fmt.Printf("Error saving key: %v\n", err)
		return
	}

	originalText := []byte("Sensitive data requiring protection")
	fmt.Printf("Original text: %s\n", originalText)

	encrypted, err := encryptData(originalText, key)
	if err != nil {
		fmt.Printf("Error encrypting: %v\n", err)
		return
	}
	fmt.Printf("Encrypted (base64): %s\n", base64.StdEncoding.EncodeToString(encrypted))

	loadedKey, err := loadKeyFromFile("encryption.key")
	if err != nil {
		fmt.Printf("Error loading key: %v\n", err)
		return
	}

	decrypted, err := decryptData(encrypted, loadedKey)
	if err != nil {
		fmt.Printf("Error decrypting: %v\n", err)
		return
	}
	fmt.Printf("Decrypted text: %s\n", decrypted)
}