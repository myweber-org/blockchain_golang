package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
    "os"
)

func encryptFile(inputPath, outputPath string, key []byte) error {
    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
        return fmt.Errorf("write output file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read input file: %w", err)
    }

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
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
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
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        fmt.Printf("Generate key failed: %v\n", err)
        return
    }

    testData := []byte("Sensitive information requiring encryption")
    if err := os.WriteFile("test_input.txt", testData, 0644); err != nil {
        fmt.Printf("Create test file failed: %v\n", err)
        return
    }

    if err := encryptFile("test_input.txt", "encrypted.dat", key); err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        return
    }
    fmt.Println("File encrypted successfully")

    if err := decryptFile("encrypted.dat", "decrypted.txt", key); err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        return
    }
    fmt.Println("File decrypted successfully")

    decrypted, _ := os.ReadFile("decrypted.txt")
    if string(decrypted) == string(testData) {
        fmt.Println("Verification: Data matches original")
    } else {
        fmt.Println("Verification: Data mismatch")
    }

    os.Remove("test_input.txt")
    os.Remove("encrypted.dat")
    os.Remove("decrypted.txt")
}