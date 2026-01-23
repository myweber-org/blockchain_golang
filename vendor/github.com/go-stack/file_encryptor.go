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

func encryptString(plaintext string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptString(ciphertext string, key []byte) (string, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    if len(data) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    iv := data[:aes.BlockSize]
    data = data[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(data, data)

    return string(data), nil
}

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    key, err := generateKey()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Key generation failed: %v\n", err)
        os.Exit(1)
    }

    original := "Sensitive data requiring protection"
    
    encrypted, err := encryptString(original, key)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Encryption failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Encrypted: %s\n", encrypted)
    
    decrypted, err := decryptString(encrypted, key)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Decryption failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Decrypted: %s\n", decrypted)
    
    if original == decrypted {
        fmt.Println("Encryption/decryption successful")
    } else {
        fmt.Println("Encryption/decryption failed")
        os.Exit(1)
    }
}