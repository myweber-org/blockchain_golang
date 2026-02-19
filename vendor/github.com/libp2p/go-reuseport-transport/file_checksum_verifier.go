package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
)

func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func verifyChecksum(filePath, expectedChecksum string) (bool, error) {
	actualChecksum, err := calculateChecksum(filePath)
	if err != nil {
		return false, err
	}
	return actualChecksum == expectedChecksum, nil
}

func main() {
	filePath := flag.String("file", "", "Path to the file")
	expectedChecksum := flag.String("checksum", "", "Expected SHA-256 checksum")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Error: File path is required")
		flag.Usage()
		os.Exit(1)
	}

	if *expectedChecksum == "" {
		checksum, err := calculateChecksum(*filePath)
		if err != nil {
			fmt.Printf("Error calculating checksum: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("SHA-256 checksum: %s\n", checksum)
	} else {
		match, err := verifyChecksum(*filePath, *expectedChecksum)
		if err != nil {
			fmt.Printf("Error verifying checksum: %v\n", err)
			os.Exit(1)
		}
		if match {
			fmt.Println("Checksum verification PASSED")
		} else {
			fmt.Println("Checksum verification FAILED")
			os.Exit(1)
		}
	}
}