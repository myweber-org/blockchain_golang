
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataCleaner struct {
    inputPath  string
    outputPath string
    seenRows   map[string]bool
}

func NewDataCleaner(input, output string) *DataCleaner {
    return &DataCleaner{
        inputPath:  input,
        outputPath: output,
        seenRows:   make(map[string]bool),
    }
}

func (dc *DataCleaner) Clean() error {
    inputFile, err := os.Open(dc.inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dc.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read record: %w", err)
        }

        rowKey := strings.Join(record, "|")
        if !dc.seenRows[rowKey] {
            dc.seenRows[rowKey] = true
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("failed to write record: %w", err)
            }
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    cleaner := NewDataCleaner(os.Args[1], os.Args[2])
    if err := cleaner.Clean(); err != nil {
        fmt.Printf("Error cleaning data: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Data cleaned successfully. Output written to %s\n", os.Args[2])
}