
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("failed to open %s: %v", path, err))
		fp.mu.Unlock()
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("error scanning %s: %v", path, err))
		fp.mu.Unlock()
		return
	}

	fp.mu.Lock()
	fp.processed++
	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	fp.mu.Unlock()
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	var wg sync.WaitGroup

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		wg.Add(1)
		go fp.ProcessFile(path, &wg)
		return nil
	})

	if err != nil {
		return fmt.Errorf("walk error: %v", err)
	}

	wg.Wait()
	return nil
}

func (fp *FileProcessor) Stats() {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("Files processed: %d\n", fp.processed)
	if len(fp.errors) > 0 {
		fmt.Printf("Errors encountered: %d\n", len(fp.errors))
		for _, err := range fp.errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	processor := NewFileProcessor()
	
	if err := processor.ProcessDirectory(os.Args[1]); err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}
	
	processor.Stats()
}