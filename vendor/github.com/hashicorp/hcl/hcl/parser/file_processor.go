package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers   int
	batchSize int
	mu        sync.RWMutex
	stats     map[string]int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(fileChan, &wg)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	return nil
}

func (fp *FileProcessor) worker(files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range files {
		if err := fp.processSingleFile(filePath); err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			continue
		}

		fp.mu.Lock()
		fp.stats["processed"]++
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) processSingleFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	reader := bufio.NewReader(file)
	lineCount := 0

	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		lineCount++
	}

	fp.mu.Lock()
	fp.stats[ext] = fp.stats[ext] + lineCount
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", path, lineCount)
	return nil
}

func (fp *FileProcessor) GetStats() map[string]int {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	statsCopy := make(map[string]int)
	for k, v := range fp.stats {
		statsCopy[k] = v
	}
	return statsCopy
}

func main() {
	processor := NewFileProcessor(4, 10)

	sampleFiles := []string{
		"data/file1.txt",
		"data/file2.log",
		"data/file3.csv",
	}

	start := time.Now()
	if err := processor.ProcessFiles(sampleFiles); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		return
	}
	elapsed := time.Since(start)

	stats := processor.GetStats()
	fmt.Printf("\nProcessing completed in %v\n", elapsed)
	fmt.Println("Statistics:")
	for key, value := range stats {
		fmt.Printf("  %s: %d\n", key, value)
	}
}