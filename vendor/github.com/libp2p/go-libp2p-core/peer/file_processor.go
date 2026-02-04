
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
	stats     ProcessingStats
}

type ProcessingStats struct {
	FilesProcessed int
	TotalBytes     int64
	Errors         int
	StartTime      time.Time
	EndTime        time.Time
}

type FileTask struct {
	Path    string
	Content []byte
	Err     error
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 4
	}
	if batchSize < 1 {
		batchSize = 100
	}
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     ProcessingStats{},
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string, pattern string) error {
	fp.stats.StartTime = time.Now()
	defer func() { fp.stats.EndTime = time.Now() }()

	matches, err := filepath.Glob(filepath.Join(dirPath, pattern))
	if err != nil {
		return fmt.Errorf("failed to match files: %w", err)
	}

	if len(matches) == 0 {
		return errors.New("no files matched the pattern")
	}

	taskChan := make(chan FileTask, len(matches))
	resultChan := make(chan FileTask, len(matches))
	var wg sync.WaitGroup

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(taskChan, resultChan, &wg)
	}

	for _, match := range matches {
		taskChan <- FileTask{Path: match}
	}
	close(taskChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fp.mu.Lock()
		if result.Err != nil {
			fp.stats.Errors++
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Err)
		} else {
			fp.stats.FilesProcessed++
			fp.stats.TotalBytes += int64(len(result.Content))
		}
		fp.mu.Unlock()
	}

	return nil
}

func (fp *FileProcessor) worker(tasks <-chan FileTask, results chan<- FileTask, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		content, err := fp.readFile(task.Path)
		task.Content = content
		task.Err = err
		results <- task
	}
}

func (fp *FileProcessor) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content []byte
	buffer := make([]byte, fp.batchSize)

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			content = append(content, buffer[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return content, fmt.Errorf("read error: %w", err)
		}
	}

	return content, nil
}

func (fp *FileProcessor) GetStats() ProcessingStats {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.stats
}

func (fp *FileProcessor) PrintStats() {
	stats := fp.GetStats()
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Println("=== Processing Statistics ===")
	fmt.Printf("Files processed: %d\n", stats.FilesProcessed)
	fmt.Printf("Total bytes: %d\n", stats.TotalBytes)
	fmt.Printf("Errors encountered: %d\n", stats.Errors)
	fmt.Printf("Processing time: %v\n", duration.Round(time.Millisecond))
	if stats.FilesProcessed > 0 && duration > 0 {
		throughput := float64(stats.TotalBytes) / duration.Seconds()
		fmt.Printf("Throughput: %.2f bytes/second\n", throughput)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory> [pattern]")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	pattern := "*"
	if len(os.Args) > 2 {
		pattern = os.Args[2]
	}

	processor := NewFileProcessor(4, 4096)
	
	fmt.Printf("Processing files in %s matching pattern: %s\n", dirPath, pattern)
	
	if err := processor.ProcessDirectory(dirPath, pattern); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}
	
	processor.PrintStats()
}