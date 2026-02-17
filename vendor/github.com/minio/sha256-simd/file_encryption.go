
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

type Encryptor struct {
    key []byte
}

func NewEncryptor(passphrase string) *Encryptor {
    hash := sha256.Sum256([]byte(passphrase))
    return &Encryptor{key: hash[:]}
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
    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryptor) Decrypt(encrypted string) (string, error) {
    data, err := base64.URLEncoding.DecodeString(encrypted)
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
    encryptor := NewEncryptor("my-secret-passphrase")
    
    original := "Sensitive data that needs protection"
    
    encrypted, err := encryptor.Encrypt(original)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }
    
    fmt.Printf("Encrypted: %s\n", encrypted)
    
    decrypted, err := encryptor.Decrypt(encrypted)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }
    
    fmt.Printf("Decrypted: %s\n", decrypted)
    
    if strings.Compare(original, decrypted) == 0 {
        fmt.Println("Encryption/decryption successful")
    } else {
        fmt.Println("Encryption/decryption failed")
    }
}