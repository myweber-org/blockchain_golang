package main

import (
	"bufio"
	"fmt"
	"os"
)

func processFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_processor.go <filename>")
		return
	}

	lines, err := processFileLines(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed %d lines:\n", len(lines))
	for i, line := range lines {
		fmt.Printf("%d: %s\n", i+1, line)
	}
}