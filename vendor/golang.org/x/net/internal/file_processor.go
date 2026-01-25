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
	mu          sync.RWMutex
	processed   map[string]bool
	workerCount int
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		processed:   make(map[string]bool),
		workerCount: workers,
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	if !filepath.IsAbs(path) {
		return errors.New("path must be absolute")
	}

	fp.mu.RLock()
	if fp.processed[path] {
		fp.mu.RUnlock()
		return fmt.Errorf("file already processed: %s", path)
	}
	fp.mu.RUnlock()

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}

	if stats.IsDir() {
		return errors.New("path is a directory, file required")
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		wordCount += countWords(line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fp.mu.Lock()
	fp.processed[path] = true
	fp.mu.Unlock()

	fmt.Printf("Processed: %s | Lines: %d | Words: %d | Size: %d bytes\n",
		filepath.Base(path), lineCount, wordCount, stats.Size())
	return nil
}

func (fp *FileProcessor) ProcessFiles(paths []string) []error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(paths))
	semaphore := make(chan struct{}, fp.workerCount)

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := fp.ProcessFile(p); err != nil {
				errChan <- err
			}
		}(path)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}
	return errors
}

func countWords(line string) int {
	inWord := false
	count := 0
	for _, r := range line {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			if !inWord {
				count++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return count
}

func main() {
	processor := NewFileProcessor(4)
	files := []string{
		"/tmp/test1.txt",
		"/tmp/test2.txt",
		"/tmp/test3.txt",
	}

	fmt.Println("Starting file processing...")
	errors := processor.ProcessFiles(files)

	if len(errors) > 0 {
		fmt.Println("\nProcessing errors:")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	} else {
		fmt.Println("\nAll files processed successfully")
	}
}