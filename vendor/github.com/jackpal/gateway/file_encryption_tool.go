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

func generateRandomKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryption_tool.go <encrypt|decrypt>")
        os.Exit(1)
    }

    operation := os.Args[1]
    key, err := generateRandomKey()
    if err != nil {
        fmt.Printf("Error generating key: %v\n", err)
        os.Exit(1)
    }

    sampleData := []byte("Sensitive information requiring encryption")

    switch operation {
    case "encrypt":
        encrypted, err := encryptData(sampleData, key)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        encoded := base64.StdEncoding.EncodeToString(encrypted)
        fmt.Printf("Encrypted data (base64): %s\n", encoded)
        fmt.Printf("Encryption key (base64): %s\n", base64.StdEncoding.EncodeToString(key))

    case "decrypt":
        if len(os.Args) < 4 {
            fmt.Println("Usage for decrypt: go run file_encryption_tool.go decrypt <base64_data> <base64_key>")
            os.Exit(1)
        }

        encryptedData, err := base64.StdEncoding.DecodeString(os.Args[2])
        if err != nil {
            fmt.Printf("Error decoding data: %v\n", err)
            os.Exit(1)
        }

        decryptionKey, err := base64.StdEncoding.DecodeString(os.Args[3])
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            os.Exit(1)
        }

        decrypted, err := decryptData(encryptedData, decryptionKey)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted data: %s\n", decrypted)

    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}