package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	daysToKeep   = 7
)

func main() {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)
	removedCount := 0

	for _, file := range files {
		if file.ModTime().Before(cutoffTime) {
			filePath := filepath.Join(tempDir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", file.Name(), err)
			} else {
				removedCount++
				fmt.Printf("Removed: %s\n", file.Name())
			}
		}
	}

	fmt.Printf("Cleanup completed. Removed %d files.\n", removedCount)
}