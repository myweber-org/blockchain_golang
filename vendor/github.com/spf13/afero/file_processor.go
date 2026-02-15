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
	errors      []string
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		errors: make([]string, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.recordError(fmt.Sprintf("failed to open %s: %v", path, err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fp.recordError(fmt.Sprintf("error scanning %s: %v", path, err))
		return
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
}

func (fp *FileProcessor) recordError(msg string) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.errors = append(fp.errors, msg)
}

func (fp *FileProcessor) Stats() (int, []string) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed, fp.errors
}

func ProcessFiles(paths []string, maxWorkers int) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	processor := NewFileProcessor()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxWorkers)

	for _, path := range paths {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(p string) {
			defer func() { <-semaphore }()
			processor.ProcessFile(p, &wg)
		}(path)
	}

	wg.Wait()
	close(semaphore)

	processed, errors := processor.Stats()
	fmt.Printf("\nProcessing complete. Files: %d, Errors: %d\n", processed, len(errors))
	if len(errors) > 0 {
		fmt.Println("Errors encountered:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	return nil
}

func main() {
	files := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
	}

	start := time.Now()
	err := ProcessFiles(files, 2)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total execution time: %v\n", elapsed)
}