package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type FileProcessor struct {
	inputPath  string
	outputPath string
	mu         sync.Mutex
}

func NewFileProcessor(input, output string) *FileProcessor {
	return &FileProcessor{
		inputPath:  input,
		outputPath: output,
	}
}

func (fp *FileProcessor) ProcessLines(transform func(string) string) error {
	if fp.inputPath == "" || fp.outputPath == "" {
		return errors.New("input or output path is empty")
	}

	inputFile, err := os.Open(fp.inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(fp.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	lineChan := make(chan string, 10)
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			lineChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("scanner error: %w", err)
		}
		close(lineChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range lineChan {
			transformed := transform(line)
			fp.mu.Lock()
			_, err := writer.WriteString(transformed + "\n")
			fp.mu.Unlock()
			if err != nil {
				errChan <- fmt.Errorf("write error: %w", err)
				return
			}
		}
	}()

	wg.Wait()
	close(errChan)

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush error: %w", err)
	}

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func main() {
	processor := NewFileProcessor("input.txt", "output.txt")
	err := processor.ProcessLines(func(line string) string {
		return "Processed: " + line
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("File processing completed successfully")
}