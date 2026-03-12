package main

import (
	"bufio"
	"fmt"
	"os"
)

func processFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Line %d: %s\n", lineNumber, line)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_processor.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	err := processFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File processing completed successfully.")
}