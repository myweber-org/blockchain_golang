package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempFilePattern = "temp_*.txt"
	maxAgeHours     = 24
)

func main() {
	dir := "."
	files, err := filepath.Glob(filepath.Join(dir, tempFilePattern))
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	now := time.Now()
	deletedCount := 0

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Error stating file %s: %v\n", file, err)
			continue
		}

		age := now.Sub(info.ModTime())
		if age.Hours() > maxAgeHours {
			err := os.Remove(file)
			if err != nil {
				fmt.Printf("Error deleting file %s: %v\n", file, err)
			} else {
				fmt.Printf("Deleted old file: %s\n", file)
				deletedCount++
			}
		}
	}

	fmt.Printf("Cleanup completed. Deleted %d files.\n", deletedCount)
}