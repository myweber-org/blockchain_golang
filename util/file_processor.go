
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
	files, err := os.ReadDir(fp.InputDir)
	if err != nil {
		return fmt.Errorf("failed to read input directory: %w", err)
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.Workers)

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go fp.worker(fileChan, &wg)
	}

	for _, file := range files {
		if !file.IsDir() {
			fileChan <- file.Name()
		}
	}

	close(fileChan)
	wg.Wait()

	return nil
}

func (fp *FileProcessor) worker(fileChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for filename := range fileChan {
		inputPath := filepath.Join(fp.InputDir, filename)
		outputPath := filepath.Join(fp.OutputDir, "processed_"+filename)

		if err := fp.processFile(inputPath, outputPath); err != nil {
			fmt.Printf("Error processing %s: %v\n", filename, err)
		}
	}
}

func (fp *FileProcessor) processFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		processedLine := fmt.Sprintf("Processed: %s\n", line)
		if _, err := writer.WriteString(processedLine); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func main() {
	processor := NewFileProcessor("./input", "./output", 4)
	if err := processor.ProcessFiles(); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("File processing completed successfully")
}