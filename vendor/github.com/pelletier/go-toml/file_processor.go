package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  map[string]int
	wg       sync.WaitGroup
	fileChan chan string
}

func NewFileProcessor(workerCount int) *FileProcessor {
	fp := &FileProcessor{
		results:  make(map[string]int),
		fileChan: make(chan string, 100),
	}

	for i := 0; i < workerCount; i++ {
		fp.wg.Add(1)
		go fp.worker()
	}

	return fp
}

func (fp *FileProcessor) worker() {
	defer fp.wg.Done()

	for filePath := range fp.fileChan {
		count, err := fp.countLines(filePath)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			continue
		}

		fp.mu.Lock()
		fp.results[filePath] = count
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount, scanner.Err()
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			fp.fileChan <- path
		}

		return nil
	})

	close(fp.fileChan)
	fp.wg.Wait()

	return err
}

func (fp *FileProcessor) GetResults() map[string]int {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	resultCopy := make(map[string]int)
	for k, v := range fp.results {
		resultCopy[k] = v
	}

	return resultCopy
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(4)

	err := processor.ProcessDirectory(dirPath)
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}

	results := processor.GetResults()
	fmt.Println("Processing results:")
	for file, count := range results {
		fmt.Printf("%s: %d lines\n", file, count)
	}
}