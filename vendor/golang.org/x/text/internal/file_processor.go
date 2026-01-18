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
    defer fp.wg.Done()

    inputPath := filepath.Join(fp.inputDir, filename)
    outputPath := filepath.Join(fp.outputDir, "processed_"+filename)

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
        fmt.Fprintln(writer, processedLine)
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }

    writer.Flush()
    return nil
}

func transformLine(line string) string {
    var result []rune
    for _, r := range line {
        if r >= 'a' && r <= 'z' {
            result = append(result, r-32)
        } else if r >= 'A' && r <= 'Z' {
            result = append(result, r+32)
        } else {
            result = append(result, r)
        }
    }
    return string(result)
}

func (fp *FileProcessor) ProcessAll(files []string) []error {
    errorChan := make(chan error, len(files))
    var errors []error

    for _, file := range files {
        fp.wg.Add(1)
        go func(f string) {
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

    errors := processor.ProcessAll(files)
    if len(errors) > 0 {
        fmt.Printf("Processing completed with %d errors\n", len(errors))
        for _, err := range errors {
            fmt.Println(err)
        }
    } else {
        fmt.Println("All files processed successfully")
    }
}