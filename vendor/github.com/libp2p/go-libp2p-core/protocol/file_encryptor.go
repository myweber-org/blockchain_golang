
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
    "path/filepath"
)

const saltSize = 16

func deriveKey(passphrase string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(passphrase))
    hash.Write(salt)
    return hash.Sum(nil)
}

func encryptFile(inputPath, outputPath, passphrase string) error {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return err
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    inputFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    if _, err := outputFile.Write(salt); err != nil {
        return err
    }

    iv := make([]byte, aes.BlockSize)
    if _, err := rand.Read(iv); err != nil {
        return err
    }

    if _, err := outputFile.Write(iv); err != nil {
        return err
    }

    stream := cipher.NewCTR(block, iv)
    writer := &cipher.StreamWriter{S: stream, W: outputFile}

    if _, err := io.Copy(writer, inputFile); err != nil {
        return err
    }

    return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(inputFile, salt); err != nil {
        return err
    }

    key := deriveKey(passphrase, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(inputFile, iv); err != nil {
        return err
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    stream := cipher.NewCTR(block, iv)
    reader := &cipher.StreamReader{S: stream, R: inputFile}

    if _, err := io.Copy(outputFile, reader); err != nil {
        return err
    }

    return nil
}

func main() {
    if len(os.Args) != 5 {
        fmt.Printf("Usage: %s <encrypt|decrypt> <input> <output> <passphrase>\n", filepath.Base(os.Args[0]))
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
        err = errors.New("invalid mode, use 'encrypt' or 'decrypt'")
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Operation completed successfully")
}