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

func (fp *FileProcessor) ProcessDirectory(dir string) error {
	var wg sync.WaitGroup
	fileChan := make(chan string, fp.workerCount)

	for i := 0; i < fp.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				if err := fp.ProcessFile(path); err != nil {
					fmt.Printf("Error processing %s: %v\n", path, err)
				}
			}
		}()
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			fileChan <- path
		}
		return nil
	})

	close(fileChan)
	wg.Wait()
	return err
}

func main() {
	processor := NewFileProcessor(4)
	start := time.Now()

	if err := processor.ProcessDirectory("."); err != nil {
		fmt.Printf("Directory processing failed: %v\n", err)
	}

	duration := time.Since(start)
	fmt.Printf("Processing completed in %v\n", duration)
}