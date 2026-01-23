
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

const (
    saltSize   = 16
    keySize    = 32
    nonceSize  = 12
)

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("cannot open input file: %w", err)
    }
    defer inputFile.Close()

    plaintext, err := io.ReadAll(inputFile)
    if err != nil {
        return fmt.Errorf("cannot read input file: %w", err)
    }

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("cannot generate salt: %w", err)
    }

    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cannot create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("cannot create GCM: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("cannot generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

    outputData := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
    outputData = append(outputData, salt...)
    outputData = append(outputData, nonce...)
    outputData = append(outputData, ciphertext...)

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("cannot create output file: %w", err)
    }
    defer outputFile.Close()

    if _, err := outputFile.Write(outputData); err != nil {
        return fmt.Errorf("cannot write output file: %w", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("cannot open input file: %w", err)
    }
    defer inputFile.Close()

    ciphertext, err := io.ReadAll(inputFile)
    if err != nil {
        return fmt.Errorf("cannot read input file: %w", err)
    }

    if len(ciphertext) < saltSize+nonceSize {
        return errors.New("file too short to be valid")
    }

    salt := ciphertext[:saltSize]
    nonce := ciphertext[saltSize : saltSize+nonceSize]
    actualCiphertext := ciphertext[saltSize+nonceSize:]

    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("cannot create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("cannot create GCM: %w", err)
    }

    plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption failed: %w", err)
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("cannot create output file: %w", err)
    }
    defer outputFile.Close()

    if _, err := outputFile.Write(plaintext); err != nil {
        return fmt.Errorf("cannot write output file: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output> <passphrase>")
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