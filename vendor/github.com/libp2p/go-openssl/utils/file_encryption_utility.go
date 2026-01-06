
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
)

func deriveKey(passphrase string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	hash.Write(salt)
	return hash.Sum(nil)
}

func encryptData(plaintext []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := deriveKey(passphrase, salt)
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
	return append(salt, ciphertext...), nil
}

func decryptData(ciphertext []byte, passphrase string) ([]byte, error) {
	if len(ciphertext) < 16 {
		return nil, errors.New("ciphertext too short")
	}

	salt := ciphertext[:16]
	ciphertext = ciphertext[16:]

	key := deriveKey(passphrase, salt)
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
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	secretMessage := "Confidential data requiring secure storage"
	passphrase := "StrongPassphrase123!"

	fmt.Println("Original:", secretMessage)

	encrypted, err := encryptData([]byte(secretMessage), passphrase)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}
	fmt.Println("Encrypted (hex):", hex.EncodeToString(encrypted))

	decrypted, err := decryptData(encrypted, passphrase)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Println("Decrypted:", string(decrypted))

	if string(decrypted) == secretMessage {
		fmt.Println("✓ Encryption/decryption successful")
	} else {
		fmt.Println("✗ Data mismatch")
		os.Exit(1)
	}
}