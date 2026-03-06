
package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateExtension checks if a file has an allowed extension.
func ValidateExtension(filename string, allowedExtensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range allowedExtensions {
		if ext == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}

// ProcessFile simulates processing a file if its extension is valid.
func ProcessFile(filename string, allowedExtensions []string) error {
	if !ValidateExtension(filename, allowedExtensions) {
		return fmt.Errorf("file %s has an invalid extension", filename)
	}
	fmt.Printf("Processing file: %s\n", filename)
	return nil
}

func main() {
	allowed := []string{".txt", ".csv", ".json"}
	files := []string{"data.txt", "report.csv", "image.png", "config.json"}

	for _, file := range files {
		err := ProcessFile(file, allowed)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}