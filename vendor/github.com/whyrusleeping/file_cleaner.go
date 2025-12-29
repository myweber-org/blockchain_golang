package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dir := "./tmp"
	err := removeOldFiles(dir, 7*24*time.Hour)
	if err != nil {
		fmt.Printf("Error cleaning directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func removeOldFiles(dir string, maxAge time.Duration) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if time.Since(info.ModTime()) > maxAge {
			fmt.Printf("Removing old file: %s\n", path)
			return os.Remove(path)
		}
		return nil
	})
}