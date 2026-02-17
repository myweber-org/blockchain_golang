
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
	results  []string
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make([]string, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	fp.wg.Add(1)
	
	go func() {
		defer fp.wg.Done()
		
		file, err := os.Open(path)
		if err != nil {
			fp.appendResult(fmt.Sprintf("Error opening %s: %v", path, err))
			return
		}
		defer file.Close()
		
		scanner := bufio.NewScanner(file)
		lineCount := 0
		
		for scanner.Scan() {
			lineCount++
		}
		
		if err := scanner.Err(); err != nil {
			fp.appendResult(fmt.Sprintf("Error scanning %s: %v", path, err))
			return
		}
		
		fp.appendResult(fmt.Sprintf("Processed %s: %d lines", filepath.Base(path), lineCount))
	}()
	
	return nil
}

func (fp *FileProcessor) appendResult(result string) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.results = append(fp.results, result)
}

func (fp *FileProcessor) WaitAndGetResults() []string {
	fp.wg.Wait()
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> [file2] ...")
		os.Exit(1)
	}
	
	processor := NewFileProcessor()
	
	for _, filePath := range os.Args[1:] {
		if err := processor.ProcessFile(filePath); err != nil {
			fmt.Printf("Failed to queue %s: %v\n", filePath, err)
		}
	}
	
	results := processor.WaitAndGetResults()
	
	fmt.Println("Processing results:")
	for _, result := range results {
		fmt.Println(result)
	}
}