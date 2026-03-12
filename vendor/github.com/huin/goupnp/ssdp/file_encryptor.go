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
		return fmt.Errorf("read file error: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation error: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %v", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	return nil
}

func main() {
	key := []byte("32-byte-long-key-here-123456789012")
	
	err := encryptFile("plain.txt", "encrypted.dat", key)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Println("File encrypted successfully")

	err = decryptFile("encrypted.dat", "decrypted.txt", key)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Println("File decrypted successfully")
}package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "os"
    "strings"
)

const (
    saltSize      = 16
    nonceSize     = 12
    keyIterations = 100000
    keyLength     = 32
)

func deriveKey(password, salt []byte) []byte {
    hash := sha256.New()
    for i := 0; i < keyIterations; i++ {
        hash.Write(append(salt, password...))
        hash.Write([]byte{byte(i & 0xFF)})
    }
    return hash.Sum(nil)[:keyLength]
}

func encryptData(plaintext, password []byte) (string, error) {
    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    key := deriveKey(password, salt)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
    combined := append(salt, append(nonce, ciphertext...)...)
    return base64.StdEncoding.EncodeToString(combined), nil
}

func decryptData(encrypted string, password []byte) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return nil, err
    }

    if len(data) < saltSize+nonceSize {
        return nil, errors.New("invalid encrypted data")
    }

    salt := data[:saltSize]
    nonce := data[saltSize : saltSize+nonceSize]
    ciphertext := data[saltSize+nonceSize:]

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
    if len(os.Args) < 4 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input_file> <output_file>")
        fmt.Println("Password will be read from environment variable ENCRYPTION_PASSWORD")
        return
    }

    operation := strings.ToLower(os.Args[1])
    inputFile := os.Args[2]
    outputFile := os.Args[3]

    password := os.Getenv("ENCRYPTION_PASSWORD")
    if password == "" {
        fmt.Println("Error: ENCRYPTION_PASSWORD environment variable not set")
        os.Exit(1)
    }

    inputData, err := os.ReadFile(inputFile)
    if err != nil {
        fmt.Printf("Error reading input file: %v\n", err)
        os.Exit(1)
    }

    switch operation {
    case "encrypt":
        encrypted, err := encryptData(inputData, []byte(password))
        if err != nil {
            fmt.Printf("Encryption error: %v\n", err)
            os.Exit(1)
        }
        if err := os.WriteFile(outputFile, []byte(encrypted), 0644); err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File encrypted successfully to %s\n", outputFile)

    case "decrypt":
        decrypted, err := decryptData(string(inputData), []byte(password))
        if err != nil {
            fmt.Printf("Decryption error: %v\n", err)
            os.Exit(1)
        }
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing output file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("File decrypted successfully to %s\n", outputFile)

    default:
        fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'")
        os.Exit(1)
    }
}