
package main

import (
	"bufio"
	"errors"
	"fmt"
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

func (fp *FileProcessor) ProcessBatch(paths []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, fp.workerCount)

	for _, path := range paths {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(p string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if err := fp.ProcessFile(p); err != nil {
				fmt.Printf("Error processing %s: %v\n", p, err)
			}
		}(path)
	}

	wg.Wait()
}

func main() {
	processor := NewFileProcessor(3)

	sampleFiles := []string{
		"/tmp/test1.txt",
		"/tmp/test2.txt",
		"/tmp/test3.txt",
		"/tmp/test4.txt",
	}

	start := time.Now()
	processor.ProcessBatch(sampleFiles)
	duration := time.Since(start)

	fmt.Printf("Batch processing completed in %v\n", duration)
}