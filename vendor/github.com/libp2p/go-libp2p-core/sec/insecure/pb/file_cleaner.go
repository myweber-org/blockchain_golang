package main

import (
    "os"
    "path/filepath"
    "time"
)

func main() {
    tempDir := os.TempDir()
    cutoff := time.Now().AddDate(0, 0, -7)

    filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info.IsDir() {
            return nil
        }
        if info.ModTime().Before(cutoff) {
            os.Remove(path)
        }
        return nil
    })
}package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_cleaner <input_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := inputFile + ".cleaned"

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	seen := make(map[string]bool)
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}
	writer.Flush()

	fmt.Printf("Duplicate lines removed. Cleaned file saved as: %s\n", outputFile)
}