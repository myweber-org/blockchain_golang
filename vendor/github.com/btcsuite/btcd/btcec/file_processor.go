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

type FileStats struct {
	Path         string
	Size         int64
	LineCount    int
	ProcessedAt  time.Time
	Error        error
}

type FileProcessor struct {
	workers    int
	results    chan FileStats
	wg         sync.WaitGroup
	mu         sync.Mutex
	totalFiles int
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workers: workers,
		results: make(chan FileStats, 100),
	}
}

func (fp *FileProcessor) ProcessDirectory(root string) ([]FileStats, error) {
	var stats []FileStats
	
	fp.wg.Add(1)
	go fp.scanDirectory(root)
	
	go func() {
		fp.wg.Wait()
		close(fp.results)
	}()
	
	for result := range fp.results {
		stats = append(stats, result)
	}
	
	return stats, nil
}

func (fp *FileProcessor) scanDirectory(path string) {
	defer fp.wg.Done()
	
	entries, err := os.ReadDir(path)
	if err != nil {
		fp.results <- FileStats{Path: path, Error: err}
		return
	}
	
	semaphore := make(chan struct{}, fp.workers)
	var dirWg sync.WaitGroup
	
	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		
		if entry.IsDir() {
			fp.wg.Add(1)
			go fp.scanDirectory(fullPath)
		} else {
			dirWg.Add(1)
			go func(p string) {
				defer dirWg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				fp.processFile(p)
			}(fullPath)
		}
	}
	
	dirWg.Wait()
}

func (fp *FileProcessor) processFile(path string) {
	stats := FileStats{
		Path:        path,
		ProcessedAt: time.Now(),
	}
	
	file, err := os.Open(path)
	if err != nil {
		stats.Error = err
		fp.results <- stats
		return
	}
	defer file.Close()
	
	info, err := file.Stat()
	if err != nil {
		stats.Error = err
		fp.results <- stats
		return
	}
	stats.Size = info.Size()
	
	lineCount, err := countLines(file)
	if err != nil {
		stats.Error = err
	} else {
		stats.LineCount = lineCount
	}
	
	fp.mu.Lock()
	fp.totalFiles++
	fp.mu.Unlock()
	
	fp.results <- stats
}

func countLines(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	lineCount := 0
	
	for scanner.Scan() {
		lineCount++
	}
	
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	
	return lineCount, nil
}

func (fp *FileProcessor) GetTotalFiles() int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.totalFiles
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}
	
	dir := os.Args[1]
	processor := NewFileProcessor(4)
	
	start := time.Now()
	stats, err := processor.ProcessDirectory(dir)
	elapsed := time.Since(start)
	
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}
	
	var totalSize int64
	var totalLines int
	errors := []error{}
	
	for _, s := range stats {
		totalSize += s.Size
		totalLines += s.LineCount
		if s.Error != nil {
			errors = append(errors, s.Error)
		}
	}
	
	fmt.Printf("Processed %d files in %v\n", processor.GetTotalFiles(), elapsed)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Errors encountered: %d\n", len(errors))
	
	if len(errors) > 0 {
		fmt.Println("\nFirst 5 errors:")
		for i := 0; i < len(errors) && i < 5; i++ {
			fmt.Printf("  %v\n", errors[i])
		}
	}
}