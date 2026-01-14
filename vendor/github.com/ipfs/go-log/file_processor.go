package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
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
		return err
	}
	defer file.Close()

	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()
		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}
		fp.mu.Lock()
		fp.results[filename] = lineCount
		fp.mu.Unlock()
	}()
	return nil
}

func (fp *FileProcessor) WaitAndDisplay() {
	fp.wg.Wait()
	fmt.Println("Processing results:")
	for filename, lines := range fp.results {
		fmt.Printf("%s: %d lines\n", filename, lines)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> [file2] ...")
		return
	}

	processor := NewFileProcessor()
	for _, filename := range os.Args[1:] {
		if err := processor.ProcessFile(filename); err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
		}
	}
	processor.WaitAndDisplay()
}