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
	wg        sync.WaitGroup
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 3
	}
	if batchSize < 1 {
		batchSize = 10
	}
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	fileChan := make(chan string, len(paths))
	resultChan := make(chan error, len(paths))

	for i := 0; i < fp.workers; i++ {
		fp.wg.Add(1)
		go fp.worker(fileChan, resultChan, processor)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	fp.wg.Wait()
	close(resultChan)

	for err := range resultChan {
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}
	}
	return nil
}

func (fp *FileProcessor) worker(files <-chan string, results chan<- error, processor func(string) error) {
	defer fp.wg.Done()
	for file := range files {
		results <- processor(file)
	}
}

func (fp *FileProcessor) BatchReadLines(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var batches [][]string
	scanner := bufio.NewScanner(file)
	currentBatch := make([]string, 0, fp.batchSize)

	for scanner.Scan() {
		currentBatch = append(currentBatch, scanner.Text())
		if len(currentBatch) >= fp.batchSize {
			batches = append(batches, currentBatch)
			currentBatch = make([]string, 0, fp.batchSize)
		}
	}

	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return batches, nil
}

func (fp *FileProcessor) CopyWithTimestamp(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	written, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	timestamp := info.ModTime().Format("20060102_150405")
	newPath := fmt.Sprintf("%s_%s%s", dst[:len(dst)-len(filepath.Ext(dst))], timestamp, filepath.Ext(dst))

	if err := os.Rename(dst, newPath); err != nil {
		return err
	}

	fmt.Printf("Copied %s to %s (%d bytes)\n", src, newPath, written)
	return nil
}

func main() {
	processor := NewFileProcessor(4, 5)

	files := []string{"input1.txt", "input2.txt"}
	err := processor.ProcessFiles(files, func(path string) error {
		fmt.Printf("Processing: %s\n", path)
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	batches, err := processor.BatchReadLines("sample.txt")
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
	} else {
		fmt.Printf("Read %d batches\n", len(batches))
	}
}