
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
    "strings"
)

type DataCleaner struct {
    InputPath  string
    OutputPath string
    Delimiter  rune
}

func NewDataCleaner(input, output string) *DataCleaner {
    return &DataCleaner{
        InputPath:  input,
        OutputPath: output,
        Delimiter:  ',',
    }
}

func (dc *DataCleaner) Clean() error {
    inputFile, err := os.Open(dc.InputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dc.OutputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    reader.Comma = dc.Delimiter
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    lineNumber := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("warning: line %d: %v", lineNumber, err)
            continue
        }

        cleaned := dc.processRecord(record)
        if cleaned != nil {
            if err := writer.Write(cleaned); err != nil {
                return fmt.Errorf("failed to write record: %w", err)
            }
        }
        lineNumber++
    }

    return nil
}

func (dc *DataCleaner) processRecord(record []string) []string {
    cleaned := make([]string, len(record))
    for i, field := range record {
        cleaned[i] = strings.TrimSpace(field)
        if cleaned[i] == "" {
            cleaned[i] = "N/A"
        }
    }
    return cleaned
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    cleaner := NewDataCleaner(os.Args[1], os.Args[2])
    if err := cleaner.Clean(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Data cleaning completed successfully")
}