package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
)

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
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, plaintext, nil), nil
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
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        fmt.Println("Key generation failed:", err)
        return
    }

    secretMessage := []byte("This is a confidential message")
    
    encrypted, err := encryptData(secretMessage, key)
    if err != nil {
        fmt.Println("Encryption failed:", err)
        return
    }
    fmt.Printf("Encrypted: %x\n", encrypted)

    decrypted, err := decryptData(encrypted, key)
    if err != nil {
        fmt.Println("Decryption failed:", err)
        return
    }
    fmt.Printf("Decrypted: %s\n", decrypted)
}