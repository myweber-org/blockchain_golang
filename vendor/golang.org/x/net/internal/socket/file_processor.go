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
	mu        sync.Mutex
	results   []ProcessResult
}

type ProcessResult struct {
	Filename string
	Size     int64
	Lines    int
	Error    error
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 1
	}
	if batchSize < 1 {
		batchSize = 10
	}
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		results:   make([]ProcessResult, 0),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) []ProcessResult {
	var wg sync.WaitGroup
	fileChan := make(chan string, len(paths))

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(fileChan, &wg)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	return fp.results
}

func (fp *FileProcessor) worker(files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	batch := make([]string, 0, fp.batchSize)
	for file := range files {
		batch = append(batch, file)

		if len(batch) >= fp.batchSize {
			fp.processBatch(batch)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		fp.processBatch(batch)
	}
}

func (fp *FileProcessor) processBatch(files []string) {
	var batchResults []ProcessResult
	for _, filename := range files {
		result := fp.analyzeFile(filename)
		batchResults = append(batchResults, result)
	}

	fp.mu.Lock()
	fp.results = append(fp.results, batchResults...)
	fp.mu.Unlock()
}

func (fp *FileProcessor) analyzeFile(filename string) ProcessResult {
	start := time.Now()
	defer func() {
		fmt.Printf("Processed %s in %v\n", filename, time.Since(start))
	}()

	file, err := os.Open(filename)
	if err != nil {
		return ProcessResult{Filename: filename, Error: err}
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return ProcessResult{Filename: filename, Error: err}
	}

	if info.IsDir() {
		return ProcessResult{Filename: filename, Error: errors.New("is a directory")}
	}

	lines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}

	if err := scanner.Err(); err != nil {
		return ProcessResult{Filename: filename, Error: err}
	}

	return ProcessResult{
		Filename: filename,
		Size:     info.Size(),
		Lines:    lines,
		Error:    nil,
	}
}

func (fp *FileProcessor) GenerateReport(w io.Writer) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	totalFiles := len(fp.results)
	var totalSize int64
	var totalLines int
	errors := 0

	for _, result := range fp.results {
		if result.Error != nil {
			errors++
			continue
		}
		totalSize += result.Size
		totalLines += result.Lines
	}

	fmt.Fprintf(w, "File Processing Report\n")
	fmt.Fprintf(w, "=====================\n")
	fmt.Fprintf(w, "Total files processed: %d\n", totalFiles)
	fmt.Fprintf(w, "Successful: %d\n", totalFiles-errors)
	fmt.Fprintf(w, "Failed: %d\n", errors)
	fmt.Fprintf(w, "Total size: %d bytes\n", totalSize)
	fmt.Fprintf(w, "Total lines: %d\n", totalLines)
	fmt.Fprintf(w, "Average file size: %.2f bytes\n", float64(totalSize)/float64(totalFiles-errors))
	fmt.Fprintf(w, "Average lines per file: %.2f\n", float64(totalLines)/float64(totalFiles-errors))

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> [file2] ...")
		os.Exit(1)
	}

	processor := NewFileProcessor(4, 5)
	results := processor.ProcessFiles(os.Args[1:])

	successCount := 0
	for _, result := range results {
		if result.Error == nil {
			successCount++
			fmt.Printf("✓ %s: %d bytes, %d lines\n",
				filepath.Base(result.Filename),
				result.Size,
				result.Lines)
		} else {
			fmt.Printf("✗ %s: %v\n",
				filepath.Base(result.Filename),
				result.Error)
		}
	}

	fmt.Printf("\nProcessed %d files (%d successful)\n", len(results), successCount)

	if err := processor.GenerateReport(os.Stdout); err != nil {
		fmt.Printf("Error generating report: %v\n", err)
	}
}