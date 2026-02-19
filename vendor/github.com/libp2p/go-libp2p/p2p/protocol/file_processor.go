package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filePath := flag.String("file", "", "Path to the input text file")
	filter := flag.String("filter", "", "Substring to filter lines (optional)")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path using -file flag")
		os.Exit(1)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if *filter == "" || strings.Contains(line, *filter) {
			fmt.Printf("%d: %s\n", lineNumber, line)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
}