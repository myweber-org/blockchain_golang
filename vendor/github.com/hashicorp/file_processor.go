package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	inputPath  string
	outputPath string
	transform  func(string) string
}

func NewFileProcessor(input, output string, transform func(string) string) *FileProcessor {
	return &FileProcessor{
		inputPath:  input,
		outputPath: output,
		transform:  transform,
	}
}

func (fp *FileProcessor) Process() error {
	if fp.transform == nil {
		return errors.New("transform function cannot be nil")
	}

	inputFile, err := os.Open(fp.inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputDir := filepath.Dir(fp.outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFile, err := os.Create(fp.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	var wg sync.WaitGroup
	lines := make(chan string, 100)
	errorsChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errorsChan <- fmt.Errorf("error reading input file: %w", err)
		}
		close(lines)
	}()

	processedLines := make(chan string, 100)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range lines {
			processedLines <- fp.transform(line)
		}
		close(processedLines)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		writer := bufio.NewWriter(outputFile)
		for processedLine := range processedLines {
			if _, err := writer.WriteString(processedLine + "\n"); err != nil {
				errorsChan <- fmt.Errorf("error writing to output file: %w", err)
				return
			}
		}
		if err := writer.Flush(); err != nil {
			errorsChan <- fmt.Errorf("error flushing output file: %w", err)
		}
	}()

	wg.Wait()
	close(errorsChan)

	if err := <-errorsChan; err != nil {
		return err
	}

	return nil
}

func main() {
	processor := NewFileProcessor(
		"input.txt",
		"output/processed.txt",
		func(s string) string {
			return "Processed: " + s
		},
	)

	if err := processor.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File processing completed successfully")
}