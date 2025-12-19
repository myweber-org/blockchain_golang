package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "strings"
)

type EncryptionResult struct {
    Ciphertext string
    Nonce      string
    Salt       string
}

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptData(plaintext, passphrase string) (*EncryptionResult, error) {
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
    if _, err := rand.Read(nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
    return &EncryptionResult{
        Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
        Nonce:      base64.StdEncoding.EncodeToString(nonce),
        Salt:       base64.StdEncoding.EncodeToString(salt),
    }, nil
}

func decryptData(result *EncryptionResult, passphrase string) (string, error) {
    salt, err := base64.StdEncoding.DecodeString(result.Salt)
    if err != nil {
        return "", err
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext, err := base64.StdEncoding.DecodeString(result.Ciphertext)
    if err != nil {
        return "", err
    }

    nonce, err := base64.StdEncoding.DecodeString(result.Nonce)
    if err != nil {
        return "", err
    }

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func main() {
    secretMessage := "Sensitive data requiring protection"
    password := "SecurePass123!"

    fmt.Println("Original message:", secretMessage)

    encrypted, err := encryptData(secretMessage, password)
    if err != nil {
        fmt.Println("Encryption error:", err)
        return
    }

    fmt.Printf("Encrypted result:\nCiphertext: %s\nNonce: %s\nSalt: %s\n",
        encrypted.Ciphertext[:30]+"...",
        encrypted.Nonce,
        encrypted.Salt)

    decrypted, err := decryptData(encrypted, password)
    if err != nil {
        fmt.Println("Decryption error:", err)
        return
    }

    fmt.Println("Decrypted message:", decrypted)

    if strings.Compare(secretMessage, decrypted) == 0 {
        fmt.Println("Encryption/decryption successful!")
    } else {
        fmt.Println("Encryption/decryption failed!")
    }
}