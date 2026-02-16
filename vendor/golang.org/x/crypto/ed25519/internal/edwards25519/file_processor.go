
package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type FileProcessor struct {
    inputDir  string
    outputDir string
    workers   int
}

func NewFileProcessor(input, output string, workers int) *FileProcessor {
    return &FileProcessor{
        inputDir:  input,
        outputDir: output,
        workers:   workers,
    }
}

func (fp *FileProcessor) ProcessFiles() error {
    files, err := os.ReadDir(fp.inputDir)
    if err != nil {
        return fmt.Errorf("failed to read input directory: %w", err)
    }

    var wg sync.WaitGroup
    fileChan := make(chan string, len(files))

    for i := 0; i < fp.workers; i++ {
        wg.Add(1)
        go fp.worker(&wg, fileChan)
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

func (fp *FileProcessor) worker(wg *sync.WaitGroup, fileChan <-chan string) {
    defer wg.Done()

    for filename := range fileChan {
        inputPath := filepath.Join(fp.inputDir, filename)
        outputPath := filepath.Join(fp.outputDir, "processed_"+filename)

        if err := fp.processFile(inputPath, outputPath); err != nil {
            fmt.Printf("error processing %s: %v\n", filename, err)
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
        processedLine := transformLine(line)
        if _, err := writer.WriteString(processedLine + "\n"); err != nil {
            return fmt.Errorf("failed to write line: %w", err)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("scanner error: %w", err)
    }

    if err := writer.Flush(); err != nil {
        return fmt.Errorf("failed to flush writer: %w", err)
    }

    return nil
}

func transformLine(line string) string {
    return "PROCESSED: " + line
}

func main() {
    processor := NewFileProcessor("./input", "./output", 4)
    if err := processor.ProcessFiles(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("File processing completed successfully")
}