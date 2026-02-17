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
	maxFileSize int64
}

func NewFileProcessor(maxSize int64) *FileProcessor {
	return &FileProcessor{
		processed:   make(map[string]bool),
		maxFileSize: maxSize,
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	if !fp.canProcess(path) {
		return errors.New("file already processed or too large")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if info.Size() > fp.maxFileSize {
		return errors.New("file exceeds size limit")
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

	fp.markProcessed(path)
	fmt.Printf("Processed %s: %d lines, %d words\n", path, lineCount, wordCount)
	return nil
}

func (fp *FileProcessor) ProcessDirectory(dir string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 10)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if err := fp.ProcessFile(p); err != nil {
				errChan <- fmt.Errorf("%s: %w", p, err)
			}
		}(path)

		return nil
	})

	wg.Wait()
	close(errChan)

	if err != nil {
		return err
	}

	for e := range errChan {
		fmt.Println("Processing error:", e)
	}

	return nil
}

func (fp *FileProcessor) canProcess(path string) bool {
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	_, exists := fp.processed[path]
	return !exists
}

func (fp *FileProcessor) markProcessed(path string) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.processed[path] = true
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
	processor := NewFileProcessor(10 * 1024 * 1024) // 10MB limit

	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	if err := processor.ProcessDirectory(dir); err != nil {
		fmt.Printf("Directory processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing completed successfully")
}