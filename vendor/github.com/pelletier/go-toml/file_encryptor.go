
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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
}

func encryptData(plaintext []byte, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encryptedText string, key []byte) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
    if err != nil {
        return nil, err
    }

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
    key, err := generateKey()
    if err != nil {
        fmt.Printf("Key generation failed: %v\n", err)
        os.Exit(1)
    }

    originalText := "Sensitive data requiring protection"
    fmt.Printf("Original text: %s\n", originalText)

    encrypted, err := encryptData([]byte(originalText), key)
    if err != nil {
        fmt.Printf("Encryption failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Encrypted text: %s\n", encrypted)

    decrypted, err := decryptData(encrypted, key)
    if err != nil {
        fmt.Printf("Decryption failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Decrypted text: %s\n", string(decrypted))
}package main

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

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return ciphertext, nil
}

func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    if len(ciphertext) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return ciphertext, nil
}

func main() {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        fmt.Println("Error generating key:", err)
        return
    }

    secretMessage := "This is a confidential message"
    fmt.Println("Original:", secretMessage)

    encrypted, err := encryptData([]byte(secretMessage), key)
    if err != nil {
        fmt.Println("Encryption error:", err)
        return
    }

    encoded := base64.StdEncoding.EncodeToString(encrypted)
    fmt.Println("Encrypted (base64):", encoded)

    decoded, _ := base64.StdEncoding.DecodeString(encoded)
    decrypted, err := decryptData(decoded, key)
    if err != nil {
        fmt.Println("Decryption error:", err)
        return
    }

    fmt.Println("Decrypted:", string(decrypted))
}