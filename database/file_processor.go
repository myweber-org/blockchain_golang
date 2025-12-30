package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	Workers    int
	BatchSize  int
	ResultChan chan ProcessResult
	ErrorChan  chan error
}

type ProcessResult struct {
	Filename string
	Lines    int
	Duration time.Duration
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		Workers:    workers,
		BatchSize:  batchSize,
		ResultChan: make(chan ProcessResult, 100),
		ErrorChan:  make(chan error, 100),
	}
}

func (fp *FileProcessor) ProcessFiles(filepaths []string) {
	var wg sync.WaitGroup
	fileBatches := fp.createBatches(filepaths)

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go fp.worker(i, fileBatches[i], &wg)
	}

	wg.Wait()
	close(fp.ResultChan)
	close(fp.ErrorChan)
}

func (fp *FileProcessor) createBatches(filepaths []string) [][]string {
	batches := make([][]string, fp.Workers)
	for i, path := range filepaths {
		batches[i%fp.Workers] = append(batches[i%fp.Workers], path)
	}
	return batches
}

func (fp *FileProcessor) worker(id int, files []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, filepath := range files {
		start := time.Now()
		lines, err := fp.countLines(filepath)
		duration := time.Since(start)

		if err != nil {
			fp.ErrorChan <- fmt.Errorf("worker %d: %v", id, err)
			continue
		}

		fp.ResultChan <- ProcessResult{
			Filename: filepath,
			Lines:    lines,
			Duration: duration,
		}
	}
}

func (fp *FileProcessor) countLines(filepath string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	if filepath.Ext(filepath) != ".txt" {
		return 0, errors.New("unsupported file format")
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func main() {
	files := []string{"data1.txt", "data2.txt", "data3.txt"}
	processor := NewFileProcessor(3, 10)

	go processor.ProcessFiles(files)

	for result := range processor.ResultChan {
		fmt.Printf("Processed %s: %d lines in %v\n",
			result.Filename, result.Lines, result.Duration)
	}

	for err := range processor.ErrorChan {
		fmt.Printf("Error: %v\n", err)
	}
}