
package main

import (
	"fmt"
	"io/fs"
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
		os.Exit(1)
	}

	dir := os.Args[1]
	err := cleanOldTempFiles(dir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Cleanup completed successfully")
}

func cleanOldTempFiles(dir string) error {
	cutoffTime := time.Now().Add(-maxAgeHours * time.Hour)

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(cutoffTime) && filepath.Base(path)[:len(tempFilePrefix)] == tempFilePrefix {
			fmt.Printf("Removing old temp file: %s\n", path)
			return os.Remove(path)
		}

		return nil
	})
}