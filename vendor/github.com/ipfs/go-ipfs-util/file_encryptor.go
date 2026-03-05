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
		return fmt.Errorf("read file error: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("cipher creation error: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("GCM mode error: %w", err)
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
		return fmt.Errorf("GCM mode error: %w", err)
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

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("key generation error: %w", err)
	}
	return key, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: file_encryptor <encrypt|decrypt> <input> <output>")
		fmt.Println("Example: file_encryptor encrypt secret.txt secret.enc")
		os.Exit(1)
	}

	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	key, err := generateKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	switch operation {
	case "encrypt":
		if err := encryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File encrypted successfully. Key: %x\n", key)
	case "decrypt":
		fmt.Print("Enter encryption key (hex): ")
		var keyHex string
		fmt.Scanln(&keyHex)
		key = []byte(keyHex)
		if err := decryptFile(inputFile, outputFile, key); err != nil {
			fmt.Printf("Decryption failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully.")
	default:
		fmt.Println("Invalid operation. Use 'encrypt' or 'decrypt'.")
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

func generateKey() ([]byte, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    return key, nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_encryptor.go <command>")
        fmt.Println("Commands: encrypt, decrypt, genkey")
        return
    }

    command := os.Args[1]

    switch command {
    case "encrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go encrypt <input_file> <key_base64>")
            return
        }

        inputFile := os.Args[2]
        keyBase64 := os.Args[3]

        key, err := base64.StdEncoding.DecodeString(keyBase64)
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            return
        }

        data, err := os.ReadFile(inputFile)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        encrypted, err := encryptData(data, key)
        if err != nil {
            fmt.Printf("Error encrypting: %v\n", err)
            return
        }

        outputFile := inputFile + ".enc"
        if err := os.WriteFile(outputFile, encrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }

        fmt.Printf("Encrypted file saved as: %s\n", outputFile)

    case "decrypt":
        if len(os.Args) != 4 {
            fmt.Println("Usage: go run file_encryptor.go decrypt <input_file> <key_base64>")
            return
        }

        inputFile := os.Args[2]
        keyBase64 := os.Args[3]

        key, err := base64.StdEncoding.DecodeString(keyBase64)
        if err != nil {
            fmt.Printf("Error decoding key: %v\n", err)
            return
        }

        data, err := os.ReadFile(inputFile)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
            return
        }

        decrypted, err := decryptData(data, key)
        if err != nil {
            fmt.Printf("Error decrypting: %v\n", err)
            return
        }

        outputFile := inputFile + ".dec"
        if err := os.WriteFile(outputFile, decrypted, 0644); err != nil {
            fmt.Printf("Error writing file: %v\n", err)
            return
        }

        fmt.Printf("Decrypted file saved as: %s\n", outputFile)

    case "genkey":
        key, err := generateKey()
        if err != nil {
            fmt.Printf("Error generating key: %v\n", err)
            return
        }

        keyBase64 := base64.StdEncoding.EncodeToString(key)
        fmt.Printf("Generated key: %s\n", keyBase64)

    default:
        fmt.Println("Unknown command. Available commands: encrypt, decrypt, genkey")
    }
}