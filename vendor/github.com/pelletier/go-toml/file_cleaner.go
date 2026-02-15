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
	maxAge       = 7 * 24 * time.Hour
	fileModePerm = 0755
)

func main() {
	if err := cleanOldFiles(); err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles() error {
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return nil
	}

	cutoffTime := time.Now().Add(-maxAge)

	return filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
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

		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
			fmt.Printf("Removed: %s\n", path)
		}
		return nil
	})
}