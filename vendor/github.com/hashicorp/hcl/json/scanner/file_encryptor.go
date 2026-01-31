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

func encrypt(plaintext []byte, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encodedCiphertext string, key []byte) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
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
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <filename>")
        os.Exit(1)
    }

    action := os.Args[1]
    filename := os.Args[2]

    data, err := os.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    key, err := generateKey()
    if err != nil {
        fmt.Printf("Error generating key: %v\n", err)
        os.Exit(1)
    }

    switch action {
    case "encrypt":
        encrypted, err := encrypt(data, key)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted content:\n%s\n", encrypted)
        fmt.Printf("Encryption key (base64): %s\n", base64.StdEncoding.EncodeToString(key))

    case "decrypt":
        fmt.Print("Enter base64-encoded key: ")
        var keyB64 string
        fmt.Scanln(&keyB64)

        key, err := base64.StdEncoding.DecodeString(keyB64)
        if err != nil {
            fmt.Printf("Invalid key format: %v\n", err)
            os.Exit(1)
        }

        decrypted, err := decrypt(string(data), key)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted content:\n%s\n", decrypted)

    default:
        fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}