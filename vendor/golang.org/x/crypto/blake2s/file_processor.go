package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	inputPath  string
	outputPath string
	mu         sync.Mutex
}

func NewFileProcessor(input, output string) *FileProcessor {
	return &FileProcessor{
		inputPath:  input,
		outputPath: output,
	}
}

func (fp *FileProcessor) ProcessLines() error {
	file, err := os.Open(fp.inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	outputFile, err := os.Create(fp.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)
	var wg sync.WaitGroup
	lineChan := make(chan string, 100)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineChan {
				processed := fp.transformLine(line)
				fp.mu.Lock()
				writer.WriteString(processed + "\n")
				fp.mu.Unlock()
			}
		}()
	}

	for scanner.Scan() {
		lineChan <- scanner.Text()
	}
	close(lineChan)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	wg.Wait()
	writer.Flush()

	return nil
}

func (fp *FileProcessor) transformLine(line string) string {
	return fmt.Sprintf("Processed: %s", line)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: file_processor <input_file> <output_file>")
		os.Exit(1)
	}

	processor := NewFileProcessor(os.Args[1], os.Args[2])
	if err := processor.ProcessLines(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File processing completed successfully")
}