package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
)

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
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
    return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    key, err := generateKey()
    if err != nil {
        fmt.Printf("Key generation error: %v\n", err)
        return
    }

    fmt.Printf("Generated key (base64): %s\n", base64.StdEncoding.EncodeToString(key))

    secretMessage := "This is a confidential message"
    encrypted, err := encrypt([]byte(secretMessage), key)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }

    fmt.Printf("Encrypted (base64): %s\n", base64.StdEncoding.EncodeToString(encrypted))

    decrypted, err := decrypt(encrypted, key)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }

    fmt.Printf("Decrypted: %s\n", string(decrypted))
}