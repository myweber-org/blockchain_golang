package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "errors"
    "fmt"
    "io"
    "os"
    "golang.org/x/crypto/pbkdf2"
)

const (
    saltSize      = 16
    keyIterations = 100000
    keyLength     = 32
)

func deriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, keyIterations, keyLength, sha256.New)
}

func encryptFile(inputPath, outputPath, password string) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("generate salt: %w", err)
    }

    key := deriveKey(password, salt)
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

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    ciphertext = append(salt, ciphertext...)

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
        return fmt.Errorf("write output file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, password string) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    if len(ciphertext) < saltSize {
        return errors.New("invalid ciphertext length")
    }

    salt := ciphertext[:saltSize]
    ciphertext = ciphertext[saltSize:]

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("invalid ciphertext length")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decrypt data: %w", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write output file: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input> <output> <password>")
        os.Exit(1)
    }

    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    password := os.Args[4]

    var err error
    switch mode {
    case "encrypt":
        err = encryptFile(inputPath, outputPath, password)
    case "decrypt":
        err = decryptFile(inputPath, outputPath, password)
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