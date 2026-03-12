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
	results  []string
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make([]string, 0),
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
		line := scanner.Text()
		fp.wg.Add(1)
		go fp.processLine(line, lineCount)
		lineCount++
	}

	fp.wg.Wait()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func (fp *FileProcessor) processLine(line string, index int) {
	defer fp.wg.Done()

	time.Sleep(10 * time.Millisecond)

	processed := fmt.Sprintf("[%d] %s", index, line)

	fp.mu.Lock()
	fp.results = append(fp.results, processed)
	fp.mu.Unlock()
}

func (fp *FileProcessor) GetResults() []string {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	processor := NewFileProcessor()

	start := time.Now()
	err := processor.ProcessFile(filename)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	results := processor.GetResults()
	fmt.Printf("Processed %d lines in %v\n", len(results), duration)

	for i, result := range results {
		if i < 5 {
			fmt.Println(result)
		}
	}

	if len(results) > 5 {
		fmt.Printf("... and %d more lines\n", len(results)-5)
	}
}