package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_cache"
	maxAgeHours  = 24
)

func main() {
	err := cleanOldFiles(tempDir, maxAgeHours)
	if err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, maxAge int) error {
	cutoffTime := time.Now().Add(-time.Duration(maxAge) * time.Hour)

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
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
			err = os.Remove(path)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				fmt.Printf("Removed: %s\n", path)
			}
		}
		return nil
	})
}