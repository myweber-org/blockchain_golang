package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  map[string]int
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		time.Sleep(1 * time.Millisecond)
	}

	fp.mu.Lock()
	fp.results[filename] = lineCount
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) ProcessConcurrently(files []string) {
	for _, file := range files {
		fp.wg.Add(1)
		go func(f string) {
			defer fp.wg.Done()
			if err := fp.ProcessFile(f); err != nil {
				fmt.Printf("Error processing %s: %v\n", f, err)
			}
		}(file)
	}
	fp.wg.Wait()
}

func (fp *FileProcessor) DisplayResults() {
	fmt.Println("Processing Results:")
	for filename, count := range fp.results {
		fmt.Printf("%s: %d lines\n", filename, count)
	}
}

func main() {
	files := []string{"data1.txt", "data2.txt", "data3.txt"}
	
	processor := NewFileProcessor()
	start := time.Now()
	
	processor.ProcessConcurrently(files)
	
	elapsed := time.Since(start)
	processor.DisplayResults()
	
	fmt.Printf("Total processing time: %v\n", elapsed)
}