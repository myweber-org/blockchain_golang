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
        _, err := writer.WriteString(processedLine + "\n")
        if err != nil {
            return fmt.Errorf("write error: %w", err)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("scan error: %w", err)
    }

    writer.Flush()
    return nil
}

func (fp *FileProcessor) ProcessAll() []error {
    entries, err := os.ReadDir(fp.inputDir)
    if err != nil {
        return []error{err}
    }

    var errors []error
    errorChan := make(chan error, len(entries))

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        fp.wg.Add(1)
        go func(filename string) {
            if err := fp.ProcessFile(filename); err != nil {
                errorChan <- err
            }
        }(entry.Name())
    }

    fp.wg.Wait()
    close(errorChan)

    for err := range errorChan {
        errors = append(errors, err)
    }

    return errors
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

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: file_processor <input_directory> <output_directory>")
        os.Exit(1)
    }

    processor := NewFileProcessor(os.Args[1], os.Args[2])
    errors := processor.ProcessAll()

    if len(errors) > 0 {
        fmt.Printf("Completed with %d errors:\n", len(errors))
        for _, err := range errors {
            fmt.Println(err)
        }
    } else {
        fmt.Println("All files processed successfully")
    }
}