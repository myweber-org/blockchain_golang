
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
    "strings"
)

const (
    saltSize   = 16
    nonceSize  = 12
    keySize    = 32
    versionTag = "v1"
)

func deriveKey(password, salt []byte) []byte {
    hash := sha256.New()
    hash.Write(password)
    hash.Write(salt)
    for i := 0; i < 10000; i++ {
        hash.Write(hash.Sum(nil))
    }
    return hash.Sum(nil)[:keySize]
}

func encryptData(plaintext, password []byte) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
    payload := append(salt, append(nonce, ciphertext...)...)
    return versionTag + ":" + hex.EncodeToString(payload), nil
}

func decryptData(encrypted string, password []byte) ([]byte, error) {
    parts := strings.Split(encrypted, ":")
    if len(parts) != 2 || parts[0] != versionTag {
        return nil, errors.New("invalid format or version")
    }

    payload, err := hex.DecodeString(parts[1])
    if err != nil {
        return nil, err
    }

    if len(payload) < saltSize+nonceSize {
        return nil, errors.New("payload too short")
    }

    salt := payload[:saltSize]
    nonce := payload[saltSize : saltSize+nonceSize]
    ciphertext := payload[saltSize+nonceSize:]

    key := deriveKey(password, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Usage: go run file_encryption_utility.go <encrypt|decrypt> <input_file> <password>")
        os.Exit(1)
    }

    mode := os.Args[1]
    filename := os.Args[2]
    password := os.Args[3]

    data, err := os.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    switch mode {
    case "encrypt":
        encrypted, err := encryptData(data, []byte(password))
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        outputFile := filename + ".enc"
        if err := os.WriteFile(outputFile, []byte(encrypted), 0600); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted data written to %s\n", outputFile)

    case "decrypt":
        decrypted, err := decryptData(string(data), []byte(password))
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        outputFile := strings.TrimSuffix(filename, ".enc") + ".dec"
        if err := os.WriteFile(outputFile, decrypted, 0600); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted data written to %s\n", outputFile)

    default:
        fmt.Println("Invalid mode. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}