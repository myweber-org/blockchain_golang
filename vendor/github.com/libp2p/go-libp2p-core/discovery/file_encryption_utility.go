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

func encryptFile(inputPath, outputPath, passphrase string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("generate salt: %w", err)
    }

    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

    outputData := append(salt, nonce...)
    outputData = append(outputData, ciphertext...)

    if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
        return fmt.Errorf("write output file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    if len(data) < 48 {
        return errors.New("invalid encrypted file format")
    }

    salt := data[:16]
    nonce := data[16:32]
    ciphertext := data[32:]

    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decrypt: %w", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write output file: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input> <output> <passphrase>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    passphrase := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, passphrase)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, passphrase)
    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}