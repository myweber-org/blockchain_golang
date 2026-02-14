package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
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

    return gcm.Seal(nonce, nonce, plaintext, nil), nil
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
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        fmt.Println("Key generation failed:", err)
        return
    }

    secretMessage := []byte("This is a confidential message")
    
    encrypted, err := encryptData(secretMessage, key)
    if err != nil {
        fmt.Println("Encryption failed:", err)
        return
    }
    fmt.Printf("Encrypted: %x\n", encrypted)

    decrypted, err := decryptData(encrypted, key)
    if err != nil {
        fmt.Println("Decryption failed:", err)
        return
    }
    fmt.Printf("Decrypted: %s\n", decrypted)
}
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

func encryptFile(inputPath, outputPath, keyHex string) error {
    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return fmt.Errorf("invalid key: %v", err)
    }
    if len(key) != 32 {
        return errors.New("key must be 32 bytes for AES-256")
    }

    plaintext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file error: %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
        return fmt.Errorf("write file error: %v", err)
    }

    return nil
}

func decryptFile(inputPath, outputPath, keyHex string) error {
    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return fmt.Errorf("invalid key: %v", err)
    }
    if len(key) != 32 {
        return errors.New("key must be 32 bytes for AES-256")
    }

    ciphertext, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("read file error: %v", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return fmt.Errorf("decryption failed: %v", err)
    }

    if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
        return fmt.Errorf("write file error: %v", err)
    }

    return nil
}

func generateRandomKey() (string, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return "", err
    }
    return hex.EncodeToString(key), nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage:")
        fmt.Println("  Generate key: file_encryptor -genkey")
        fmt.Println("  Encrypt: file_encryptor -encrypt input.txt output.enc key")
        fmt.Println("  Decrypt: file_encryptor -decrypt input.enc output.txt key")
        return
    }

    switch os.Args[1] {
    case "-genkey":
        key, err := generateRandomKey()
        if err != nil {
            fmt.Printf("Key generation failed: %v\n", err)
            return
        }
        fmt.Printf("Generated key: %s\n", key)

    case "-encrypt":
        if len(os.Args) != 5 {
            fmt.Println("Usage: file_encryptor -encrypt input output key")
            return
        }
        if err := encryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
        } else {
            fmt.Println("Encryption successful")
        }

    case "-decrypt":
        if len(os.Args) != 5 {
            fmt.Println("Usage: file_encryptor -decrypt input output key")
            return
        }
        if err := decryptFile(os.Args[2], os.Args[3], os.Args[4]); err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
        } else {
            fmt.Println("Decryption successful")
        }

    default:
        fmt.Println("Invalid command")
    }
}