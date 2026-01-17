
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	maxAge       = 7 * 24 * time.Hour
	filePattern  = "*.tmp"
)

func main() {
	err := cleanOldFiles(tempDir, filePattern, maxAge)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Temporary files cleaned successfully")
}

func cleanOldFiles(dir, pattern string, maxAge time.Duration) error {
	cutoffTime := time.Now().Add(-maxAge)

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, d.Name())
		if err != nil || !matched {
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
				fmt.Printf("Removed old file: %s\n", path)
			}
		}

		return nil
	})
}