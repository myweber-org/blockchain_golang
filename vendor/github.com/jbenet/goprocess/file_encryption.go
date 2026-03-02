
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "fmt"
    "io"
    "os"
)

type Encryptor struct {
    key []byte
}

func NewEncryptor(key string) (*Encryptor, error) {
    if len(key) != 64 {
        return nil, errors.New("key must be 64 hex characters for AES-256")
    }
    
    decodedKey, err := hex.DecodeString(key)
    if err != nil {
        return nil, fmt.Errorf("invalid hex key: %w", err)
    }
    
    return &Encryptor{key: decodedKey}, nil
}

func (e *Encryptor) EncryptData(plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, fmt.Errorf("cipher creation failed: %w", err)
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("GCM mode failed: %w", err)
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("nonce generation failed: %w", err)
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

func (e *Encryptor) DecryptData(ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, fmt.Errorf("cipher creation failed: %w", err)
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("GCM mode failed: %w", err)
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %w", err)
    }
    
    return plaintext, nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run file_encryption.go <encrypt|decrypt> <input_file> [output_file]")
        os.Exit(1)
    }
    
    operation := os.Args[1]
    inputFile := os.Args[2]
    outputFile := "output.bin"
    
    if len(os.Args) > 3 {
        outputFile = os.Args[3]
    }
    
    key := os.Getenv("ENCRYPTION_KEY")
    if key == "" {
        fmt.Println("ENCRYPTION_KEY environment variable not set")
        os.Exit(1)
    }
    
    encryptor, err := NewEncryptor(key)
    if err != nil {
        fmt.Printf("Encryptor initialization failed: %v\n", err)
        os.Exit(1)
    }
    
    inputData, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("File read failed: %v\n", err)
        os.Exit(1)
    }
    
    var result []byte
    switch operation {
    case "encrypt":
        result, err = encryptor.EncryptData(inputData)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted %d bytes\n", len(result))
        
    case "decrypt":
        result, err = encryptor.DecryptData(inputData)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted %d bytes\n", len(result))
        
    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
    
    err = os.WriteFile(outputFile, result, 0644)
    if err != nil {
        fmt.Printf("File write failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Operation completed. Output written to %s\n", outputFile)
}