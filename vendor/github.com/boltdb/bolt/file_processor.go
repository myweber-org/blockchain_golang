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
	workerCount int
	results     chan string
	errors      chan error
	wg          sync.WaitGroup
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workerCount: workers,
		results:     make(chan string, 100),
		errors:      make(chan error, 100),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file not found: %s", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		// Simulate processing delay
		time.Sleep(10 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fp.results <- fmt.Sprintf("Processed %s: %d lines", filepath.Base(path), lineCount)
	return nil
}

func (fp *FileProcessor) worker(fileQueue <-chan string) {
	defer fp.wg.Done()
	for file := range fileQueue {
		if err := fp.ProcessFile(file); err != nil {
			fp.errors <- err
		}
	}
}

func (fp *FileProcessor) Run(files []string) {
	fileQueue := make(chan string, len(files))
	for _, file := range files {
		fileQueue <- file
	}
	close(fileQueue)

	for i := 0; i < fp.workerCount; i++ {
		fp.wg.Add(1)
		go fp.worker(fileQueue)
	}

	fp.wg.Wait()
	close(fp.results)
	close(fp.errors)
}

func (fp *FileProcessor) GetResults() []string {
	var output []string
	for result := range fp.results {
		output = append(output, result)
	}
	return output
}

func (fp *FileProcessor) GetErrors() []error {
	var errs []error
	for err := range fp.errors {
		errs = append(errs, err)
	}
	return errs
}

func main() {
	files := []string{"data1.txt", "data2.txt", "data3.txt"}
	processor := NewFileProcessor(3)

	go func() {
		processor.Run(files)
	}()

	// Collect results
	for _, result := range processor.GetResults() {
		fmt.Println(result)
	}

	// Report errors
	for _, err := range processor.GetErrors() {
		fmt.Printf("Error: %v\n", err)
	}
}