package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func computeChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

func verifyChecksum(filePath, expectedChecksum string) (bool, error) {
	actualChecksum, err := computeChecksum(filePath)
	if err != nil {
		return false, err
	}
	return actualChecksum == expectedChecksum, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run file_checksum_verifier.go <file_path> <expected_checksum>")
		fmt.Println("Or: go run file_checksum_verifier.go compute <file_path>")
		os.Exit(1)
	}

	command := os.Args[1]
	filePath := os.Args[2]

	switch command {
	case "compute":
		checksum, err := computeChecksum(filePath)
		if err != nil {
			fmt.Printf("Error computing checksum: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("SHA256 checksum for %s: %s\n", filePath, checksum)

	default:
		expectedChecksum := command
		match, err := verifyChecksum(filePath, expectedChecksum)
		if err != nil {
			fmt.Printf("Error verifying checksum: %v\n", err)
			os.Exit(1)
		}
		if match {
			fmt.Println("Checksum verification passed.")
		} else {
			fmt.Println("Checksum verification failed.")
		}
	}
}