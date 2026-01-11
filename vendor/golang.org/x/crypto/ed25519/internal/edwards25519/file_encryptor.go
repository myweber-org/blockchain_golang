
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
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
		return fmt.Errorf("GCM creation error: %v", err)
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
		return fmt.Errorf("GCM creation error: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
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

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("key generation error: %v", err)
	}
	return key, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt> <input> <output>")
		fmt.Println("Example: go run file_encryptor.go encrypt secret.txt secret.enc")
		os.Exit(1)
	}

	action := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	key, err := generateRandomKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	switch action {
	case "encrypt":
		if err := encryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted successfully.\nKey: %s\n", hex.EncodeToString(key))
	case "decrypt":
		fmt.Print("Enter encryption key (hex): ")
		var keyHex string
		fmt.Scanln(&keyHex)
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			fmt.Printf("Invalid key format: %v\n", err)
			os.Exit(1)
		}
		if err := decryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully.")
	default:
		fmt.Println("Invalid action. Use 'encrypt' or 'decrypt'.")
		os.Exit(1)
	}
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
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptData(encodedCiphertext string, key []byte) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryptor.go <encrypt|decrypt|keygen>")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "keygen":
        key, err := generateKey()
        if err != nil {
            fmt.Printf("Key generation failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Generated key: %s\n", base64.StdEncoding.EncodeToString(key))

    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go encrypt <key> <plaintext>")
            os.Exit(1)
        }
        key, err := base64.StdEncoding.DecodeString(os.Args[2])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            os.Exit(1)
        }
        encrypted, err := encryptData([]byte(os.Args[3]), key)
        if err != nil {
            fmt.Printf("Encryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Encrypted: %s\n", encrypted)

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go decrypt <key> <ciphertext>")
            os.Exit(1)
        }
        key, err := base64.StdEncoding.DecodeString(os.Args[2])
        if err != nil {
            fmt.Printf("Invalid key: %v\n", err)
            os.Exit(1)
        }
        decrypted, err := decryptData(os.Args[3], key)
        if err != nil {
            fmt.Printf("Decryption failed: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Decrypted: %s\n", string(decrypted))

    default:
        fmt.Println("Unknown command")
        os.Exit(1)
    }
}