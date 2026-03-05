
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int
	Content   string
	Valid     bool
	Timestamp time.Time
}

type Processor struct {
	records []DataRecord
	mu      sync.RWMutex
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
	}
}

func (p *Processor) AddRecord(content string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        len(p.records) + 1,
		Content:   content,
		Valid:     len(content) > 0,
		Timestamp: time.Now(),
	}

	p.records = append(p.records, record)
}

func (p *Processor) ValidateRecords() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.records) == 0 {
		return errors.New("no records to validate")
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(p.records))

	for i := range p.records {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			record := &p.records[idx]
			
			if !record.Valid {
				errorChan <- fmt.Errorf("record %d failed validation", record.ID)
				return
			}
			
			if len(record.Content) > 100 {
				record.Valid = false
				errorChan <- fmt.Errorf("record %d content too long", record.ID)
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	var validationErrors []error
	for err := range errorChan {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed with %d errors", len(validationErrors))
	}

	return nil
}

func (p *Processor) GetStats() (int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	validCount := 0
	for _, record := range p.records {
		if record.Valid {
			validCount++
		}
	}

	return len(p.records), validCount
}

func main() {
	processor := NewProcessor()

	sampleData := []string{
		"Sample data record one",
		"",
		"Another valid data entry",
		"This record has extremely long content that should fail validation because it exceeds the maximum allowed length of 100 characters",
		"Short",
	}

	for _, data := range sampleData {
		processor.AddRecord(data)
	}

	total, valid := processor.GetStats()
	fmt.Printf("Initial stats: %d total records, %d valid\n", total, valid)

	if err := processor.ValidateRecords(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}

	total, valid = processor.GetStats()
	fmt.Printf("Final stats: %d total records, %d valid\n", total, valid)
}package main

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
	Lines    int
	Size     int64
	Duration time.Duration
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

	for file := range files {
		result := fp.processSingleFile(file)
		fp.mu.Lock()
		fp.results = append(fp.results, result)
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) processSingleFile(path string) ProcessResult {
	start := time.Now()
	info, err := os.Stat(path)
	if err != nil {
		return ProcessResult{
			Filename: path,
			Error:    err,
			Duration: time.Since(start),
		}
	}

	if info.IsDir() {
		return ProcessResult{
			Filename: path,
			Error:    errors.New("path is directory"),
			Duration: time.Since(start),
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return ProcessResult{
			Filename: path,
			Error:    err,
			Duration: time.Since(start),
		}
	}
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return ProcessResult{
			Filename: path,
			Lines:    lineCount,
			Size:     info.Size(),
			Error:    err,
			Duration: time.Since(start),
		}
	}

	return ProcessResult{
		Filename: path,
		Lines:    lineCount,
		Size:     info.Size(),
		Duration: time.Since(start),
	}
}

func findFilesByPattern(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, errors.New("no files match pattern")
	}
	return matches, nil
}

func main() {
	files, err := findFilesByPattern("*.txt")
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		return
	}

	processor := NewFileProcessor(4, 20)
	results := processor.ProcessFiles(files)

	totalLines := 0
	var totalSize int64 = 0
	var totalDuration time.Duration

	fmt.Println("Processing Results:")
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("  %s: ERROR - %v\n", result.Filename, result.Error)
			continue
		}
		fmt.Printf("  %s: %d lines, %d bytes, %v\n",
			result.Filename, result.Lines, result.Size, result.Duration)
		totalLines += result.Lines
		totalSize += result.Size
		totalDuration += result.Duration
	}

	fmt.Printf("\nSummary: %d files, %d total lines, %d total bytes, %v total time\n",
		len(results), totalLines, totalSize, totalDuration)
}