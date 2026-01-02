package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	InputDir  string
	OutputDir string
	Workers   int
}

func NewFileProcessor(inputDir, outputDir string, workers int) *FileProcessor {
	return &FileProcessor{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Workers:   workers,
	}
}

func (fp *FileProcessor) ProcessFiles() error {
	files, err := filepath.Glob(filepath.Join(fp.InputDir, "*.txt"))
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	jobs := make(chan string, len(files))
	results := make(chan error, len(files))
	var wg sync.WaitGroup

	for w := 0; w < fp.Workers; w++ {
		wg.Add(1)
		go fp.worker(w, jobs, results, &wg)
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}
	}

	return nil
}

func (fp *FileProcessor) worker(id int, jobs <-chan string, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range jobs {
		start := time.Now()
		err := fp.processSingleFile(file)
		elapsed := time.Since(start)

		if err != nil {
			results <- fmt.Errorf("worker %d failed on %s: %w", id, file, err)
		} else {
			fmt.Printf("Worker %d processed %s in %v\n", id, filepath.Base(file), elapsed)
			results <- nil
		}
	}
}

func (fp *FileProcessor) processSingleFile(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	outputPath := filepath.Join(fp.OutputDir, filepath.Base(inputPath))
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		processedLine := fmt.Sprintf("Processed: %s\n", line)
		if _, err := writer.WriteString(processedLine); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush error: %w", err)
	}

	fmt.Printf("Processed %d lines from %s\n", lineCount, filepath.Base(inputPath))
	return nil
}

func main() {
	processor := NewFileProcessor("./input", "./output", 4)

	if err := os.MkdirAll("./input", 0755); err != nil {
		fmt.Printf("Error creating input directory: %v\n", err)
		return
	}

	if err := os.MkdirAll("./output", 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	fmt.Println("Starting file processing...")
	if err := processor.ProcessFiles(); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File processing completed successfully")
}