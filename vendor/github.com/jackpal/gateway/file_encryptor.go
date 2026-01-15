
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

const (
    saltSize      = 16
    keySize       = 32
    nonceSize     = 12
    versionHeader = "ENCv1"
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
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    salt := make([]byte, saltSize)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("failed to generate salt: %w", err)
    }

    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return fmt.Errorf("failed to generate nonce: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    if _, err := outputFile.Write([]byte(versionHeader)); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }
    if _, err := outputFile.Write(salt); err != nil {
        return fmt.Errorf("failed to write salt: %w", err)
    }
    if _, err := outputFile.Write(nonce); err != nil {
        return fmt.Errorf("failed to write nonce: %w", err)
    }

    buf := make([]byte, 4096)
    for {
        n, err := inputFile.Read(buf)
        if err != nil && err != io.EOF {
            return fmt.Errorf("failed to read input: %w", err)
        }
        if n == 0 {
            break
        }

        ciphertext := gcm.Seal(nil, nonce, buf[:n], nil)
        if _, err := outputFile.Write(ciphertext); err != nil {
            return fmt.Errorf("failed to write encrypted data: %w", err)
        }
    }

    return nil
}

func decryptFile(inputPath, outputPath, passphrase string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    header := make([]byte, len(versionHeader))
    if _, err := io.ReadFull(inputFile, header); err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }
    if string(header) != versionHeader {
        return errors.New("invalid file format or version")
    }

    salt := make([]byte, saltSize)
    if _, err := io.ReadFull(inputFile, salt); err != nil {
        return fmt.Errorf("failed to read salt: %w", err)
    }

    nonce := make([]byte, nonceSize)
    if _, err := io.ReadFull(inputFile, nonce); err != nil {
        return fmt.Errorf("failed to read nonce: %w", err)
    }

    key := deriveKey(passphrase, salt)

    block, err := aes.NewCipher(key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    buf := make([]byte, 4096+gcm.Overhead())
    for {
        n, err := inputFile.Read(buf)
        if err != nil && err != io.EOF {
            return fmt.Errorf("failed to read input: %w", err)
        }
        if n == 0 {
            break
        }

        plaintext, err := gcm.Open(nil, nonce, buf[:n], nil)
        if err != nil {
            return fmt.Errorf("failed to decrypt data: %w", err)
        }

        if _, err := outputFile.Write(plaintext); err != nil {
            return fmt.Errorf("failed to write decrypted data: %w", err)
        }
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
        fmt.Printf("Invalid mode: %s\n", mode)
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Operation completed successfully: %s -> %s\n", inputPath, outputPath)
}package main

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
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm creation error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce generation error: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath string, key []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("gcm creation error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption error: %w", err)
	}

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func main() {
	key := []byte("examplekey1234567890123456789012") // 32 bytes for AES-256

	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
		return
	}

	action := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	switch action {
	case "encrypt":
		if err := encryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
		} else {
			fmt.Println("Encryption successful")
		}
	case "decrypt":
		if err := decryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
		} else {
			fmt.Println("Decryption successful")
		}
	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'")
	}
}