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
	maxWorkers  int
	results     chan string
	errors      chan error
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		processed:  make(map[string]bool),
		maxWorkers: workers,
		results:    make(chan string, 100),
		errors:     make(chan error, 100),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	fp.mu.Lock()
	if fp.processed[path] {
		fp.mu.Unlock()
		return errors.New("file already processed")
	}
	fp.processed[path] = true
	fp.mu.Unlock()

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		_ = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	fp.results <- fmt.Sprintf("Processed %s: %d lines", filepath.Base(path), lineCount)
	return nil
}

func (fp *FileProcessor) ProcessDirectory(dir string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, fp.maxWorkers)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fp.errors <- err
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".txt" && filepath.Ext(path) != ".log" {
			return nil
		}

		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := fp.ProcessFile(p); err != nil {
				fp.errors <- err
			}
		}(path)

		return nil
	})

	if err != nil {
		fp.errors <- err
	}

	wg.Wait()
	close(fp.results)
	close(fp.errors)
}

func (fp *FileProcessor) Run(dir string) {
	go fp.ProcessDirectory(dir)

	done := make(chan bool)
	go func() {
		for result := range fp.results {
			fmt.Println(result)
		}
		done <- true
	}()

	go func() {
		for err := range fp.errors {
			fmt.Printf("Error: %v\n", err)
		}
		done <- true
	}()

	<-done
	<-done
	time.Sleep(100 * time.Millisecond)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	processor := NewFileProcessor(5)
	processor.Run(os.Args[1])
}