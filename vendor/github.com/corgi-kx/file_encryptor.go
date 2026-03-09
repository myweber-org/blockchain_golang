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

    key, err := generateRandomKey()
    if err != nil {
        fmt.Printf("Error generating key: %v\n", err)
        os.Exit(1)
    }

    switch action {
    case "encrypt":
        encrypted, err := encryptData(data, key)
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        encryptedFile := filename + ".enc"
        if err := os.WriteFile(encryptedFile, encrypted, 0644); err != nil {
            fmt.Printf("Error writing encrypted file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File encrypted successfully: %s\n", encryptedFile)
        fmt.Printf("Encryption key (base64): %s\n", base64.StdEncoding.EncodeToString(key))

    case "decrypt":
        var keyInput string
        fmt.Print("Enter decryption key (base64): ")
        fmt.Scanln(&keyInput)

        key, err := base64.StdEncoding.DecodeString(keyInput)
        if err != nil {
            fmt.Printf("Invalid key format: %v\n", err)
            os.Exit(1)
        }

        decrypted, err := decryptData(data, key)
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }

        decryptedFile := filename + ".dec"
        if err := os.WriteFile(decryptedFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing decrypted file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File decrypted successfully: %s\n", decryptedFile)

    default:
        fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}