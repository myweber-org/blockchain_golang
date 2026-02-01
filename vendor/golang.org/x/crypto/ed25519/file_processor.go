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
	mu          sync.Mutex
	processed   int
	errors      []error
	concurrency int
}

func NewFileProcessor(workers int) *FileProcessor {
	if workers < 1 {
		workers = 4
	}
	return &FileProcessor{
		concurrency: workers,
		errors:      make([]error, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		_ = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error for %s: %w", path, err)
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	return nil
}

func (fp *FileProcessor) ProcessFiles(paths []string) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, fp.concurrency)

	for _, path := range paths {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if err := fp.ProcessFile(p); err != nil {
				fp.mu.Lock()
				fp.errors = append(fp.errors, err)
				fp.mu.Unlock()
			}
		}(path)
	}

	wg.Wait()

	if len(fp.errors) > 0 {
		return fmt.Errorf("encountered %d errors during processing", len(fp.errors))
	}
	return nil
}

func (fp *FileProcessor) Stats() (int, []error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed, fp.errors
}

func findTextFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	start := time.Now()
	processor := NewFileProcessor(8)

	files, err := findTextFiles(os.Args[1])
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No text files found")
		os.Exit(0)
	}

	fmt.Printf("Found %d text files\n", len(files))

	if err := processor.ProcessFiles(files); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
	}

	processed, errors := processor.Stats()
	fmt.Printf("\nProcessing completed in %v\n", time.Since(start))
	fmt.Printf("Successfully processed: %d files\n", processed)
	fmt.Printf("Errors encountered: %d\n", len(errors))

	if len(errors) > 0 {
		fmt.Println("\nError details:")
		for _, e := range errors {
			fmt.Printf("  - %v\n", e)
		}
	}
}