package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoff := time.Now().AddDate(0, 0, -7)
	var removedCount int

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err == nil {
				removedCount++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	fmt.Printf("Cleaned %d temporary files older than 7 days\n", removedCount)
}package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const retentionDays = 7

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s\n", path)
			os.Remove(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}