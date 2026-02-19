package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempFilePrefix = "temp_"
	maxAgeHours    = 24
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_cleaner.go <directory>")
		return
	}

	dir := os.Args[1]
	err := cleanTempFiles(dir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
	}
}

func cleanTempFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if shouldRemoveFile(info) {
			fmt.Printf("Removing: %s\n", path)
			return os.Remove(path)
		}

		return nil
	})
}

func shouldRemoveFile(info os.FileInfo) bool {
	if !isTempFile(info.Name()) {
		return false
	}

	age := time.Since(info.ModTime())
	return age.Hours() > maxAgeHours
}

func isTempFile(filename string) bool {
	return len(filename) > len(tempFilePrefix) && filename[:len(tempFilePrefix)] == tempFilePrefix
}