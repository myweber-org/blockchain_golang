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

const (
    saltSize        = 16
    nonceSize       = 12
    keyIterations   = 100000
    keyLength       = 32
)

func deriveKey(password, salt []byte) []byte {
    hash := sha256.New()
    hash.Write(password)
    hash.Write(salt)
    for i := 0; i < keyIterations; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keyLength]
}

func encrypt(plaintext, password string) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    key := deriveKey([]byte(password), salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
    combined := append(salt, nonce...)
    combined = append(combined, ciphertext...)

    return base64.URLEncoding.EncodeToString(combined), nil
}

func decrypt(encrypted, password string) (string, error) {
    data, err := base64.URLEncoding.DecodeString(encrypted)
    if err != nil {
        return "", err
    }

    if len(data) < saltSize+nonceSize {
        return "", errors.New("invalid encrypted data")
    }

    salt := data[:saltSize]
    nonce := data[saltSize : saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

    key := deriveKey([]byte(password), salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func main() {
    password := "securePass123!"
    original := "Sensitive data requiring protection"

    encrypted, err := encrypt(original, password)
    if err != nil {
        fmt.Printf("Encryption error: %v\n", err)
        return
    }

    fmt.Printf("Encrypted: %s\n", strings.TrimSpace(encrypted))

    decrypted, err := decrypt(encrypted, password)
    if err != nil {
        fmt.Printf("Decryption error: %v\n", err)
        return
    }

    fmt.Printf("Decrypted: %s\n", decrypted)
    fmt.Printf("Match: %v\n", original == decrypted)
}