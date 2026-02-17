package main

import (
    "fmt"
    "io/ioutil"
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
    files, err := ioutil.ReadDir(fp.inputDir)
    if err != nil {
        return fmt.Errorf("failed to read input directory: %w", err)
    }

    jobs := make(chan string, len(files))
    results := make(chan error, len(files))
    var wg sync.WaitGroup

    for w := 0; w < fp.workers; w++ {
        wg.Add(1)
        go fp.worker(jobs, results, &wg)
    }

    for _, file := range files {
        if !file.IsDir() {
            jobs <- file.Name()
        }
    }
    close(jobs)

    wg.Wait()
    close(results)

    for err := range results {
        if err != nil {
            return err
        }
    }

    return nil
}

func (fp *FileProcessor) worker(jobs <-chan string, results chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()

    for filename := range jobs {
        inputPath := filepath.Join(fp.inputDir, filename)
        outputPath := filepath.Join(fp.outputDir, filename)

        data, err := ioutil.ReadFile(inputPath)
        if err != nil {
            results <- fmt.Errorf("failed to read file %s: %w", filename, err)
            continue
        }

        processedData := processContent(data)

        if err := os.MkdirAll(fp.outputDir, 0755); err != nil {
            results <- fmt.Errorf("failed to create output directory: %w", err)
            continue
        }

        if err := ioutil.WriteFile(outputPath, processedData, 0644); err != nil {
            results <- fmt.Errorf("failed to write file %s: %w", filename, err)
            continue
        }

        results <- nil
    }
}

func processContent(data []byte) []byte {
    processed := make([]byte, len(data))
    for i, b := range data {
        processed[i] = b ^ 0xFF
    }
    return processed
}

func main() {
    processor := NewFileProcessor("./input", "./output", 4)
    if err := processor.ProcessFiles(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("File processing completed successfully")
}