
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
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
		return errors.New("file already processed")
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
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fp.mu.Lock()
	fp.processed[path] = true
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	return nil
}

func (fp *FileProcessor) ProcessBatch(paths []string) []error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(paths))
	semaphore := make(chan struct{}, fp.workerCount)

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := fp.ProcessFile(p); err != nil {
				errCh <- fmt.Errorf("%s: %w", p, err)
			}
		}(path)
	}

	wg.Wait()
	close(errCh)

	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}
	return errors
}

func (fp *FileProcessor) ExportProcessedList(writer io.Writer) error {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	for path := range fp.processed {
		if _, err := fmt.Fprintln(writer, path); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	processor := NewFileProcessor(4)

	sampleFiles := []string{
		"/tmp/test1.txt",
		"/tmp/test2.txt",
	}

	for i, file := range sampleFiles {
		if err := os.WriteFile(file, []byte(fmt.Sprintf("Line %d\n", i+1)), 0644); err != nil {
			fmt.Printf("Failed to create test file: %v\n", err)
			return
		}
		defer os.Remove(file)
	}

	start := time.Now()
	errors := processor.ProcessBatch(sampleFiles)
	elapsed := time.Since(start)

	if len(errors) > 0 {
		fmt.Println("Processing errors:")
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	fmt.Printf("Batch processing completed in %v\n", elapsed)

	fmt.Println("\nProcessed files:")
	if err := processor.ExportProcessedList(os.Stdout); err != nil {
		fmt.Printf("Failed to export list: %v\n", err)
	}
}