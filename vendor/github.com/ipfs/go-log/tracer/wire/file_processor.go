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
	mu         sync.RWMutex
}

func NewFileProcessor(input, output string) *FileProcessor {
	return &FileProcessor{
		inputPath:  input,
		outputPath: output,
	}
}

func (fp *FileProcessor) ValidatePaths() error {
	if fp.inputPath == "" {
		return errors.New("input path cannot be empty")
	}
	if fp.outputPath == "" {
		return errors.New("output path cannot be empty")
	}
	return nil
}

func (fp *FileProcessor) ProcessFile() error {
	if err := fp.ValidatePaths(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	inputFile, err := os.Open(fp.inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer inputFile.Close()

	outputDir := filepath.Dir(fp.outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory: %w", err)
	}

	outputFile, err := os.Create(fp.outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer outputFile.Close()

	var wg sync.WaitGroup
	lines := make(chan string, 100)
	errors := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fp.readLines(inputFile, lines); err != nil {
			errors <- fmt.Errorf("read error: %w", err)
		}
		close(lines)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fp.writeLines(outputFile, lines); err != nil {
			errors <- fmt.Errorf("write error: %w", err)
		}
	}()

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (fp *FileProcessor) readLines(reader io.Reader, lines chan<- string) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		fp.mu.RLock()
		lines <- line
		fp.mu.RUnlock()
	}
	return scanner.Err()
}

func (fp *FileProcessor) writeLines(writer io.Writer, lines <-chan string) error {
	for line := range lines {
		fp.mu.Lock()
		_, err := fmt.Fprintln(writer, line)
		fp.mu.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	processor := NewFileProcessor("input.txt", "output/processed.txt")
	if err := processor.ProcessFile(); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("File processing completed successfully")
}