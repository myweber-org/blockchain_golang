
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
			return
		}
		defer file.Close()
		
		scanner := bufio.NewScanner(file)
		lineCount := 0
		
		for scanner.Scan() {
			lineCount++
		}
		
		if err := scanner.Err(); err != nil {
			return
		}
		
		fp.mu.Lock()
		fp.results = append(fp.results, fmt.Sprintf("%s: %d lines", filepath.Base(path), lineCount))
		fp.mu.Unlock()
	}()
	
	return nil
}

func (fp *FileProcessor) GetResults() []string {
	fp.wg.Wait()
	return fp.results
}

func main() {
	processor := NewFileProcessor()
	
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			processor.ProcessFile(file)
		}
	}
	
	results := processor.GetResults()
	
	for _, result := range results {
		fmt.Println(result)
	}
}