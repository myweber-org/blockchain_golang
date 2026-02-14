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
    wg        sync.WaitGroup
}

func NewFileProcessor(input, output string) *FileProcessor {
    return &FileProcessor{
        inputDir:  input,
        outputDir: output,
    }
}

func (fp *FileProcessor) ProcessFile(filename string) error {
    inputPath := filepath.Join(fp.inputDir, filename)
    outputPath := filepath.Join(fp.outputDir, filename+".processed")

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
        processedLine := fmt.Sprintf("PROCESSED: %s\n", line)
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

func (fp *FileProcessor) ProcessConcurrently(files []string) []error {
    errorChan := make(chan error, len(files))
    var errors []error

    for _, file := range files {
        fp.wg.Add(1)
        go func(f string) {
            defer fp.wg.Done()
            if err := fp.ProcessFile(f); err != nil {
                errorChan <- err
            }
        }(file)
    }

    fp.wg.Wait()
    close(errorChan)

    for err := range errorChan {
        errors = append(errors, err)
    }

    return errors
}

func main() {
    processor := NewFileProcessor("./input", "./output")
    files := []string{"data1.txt", "data2.txt", "data3.txt"}

    fmt.Println("Starting concurrent file processing...")
    errors := processor.ProcessConcurrently(files)

    if len(errors) > 0 {
        fmt.Printf("Encountered %d errors:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("  - %v\n", err)
        }
    } else {
        fmt.Println("All files processed successfully")
    }
}