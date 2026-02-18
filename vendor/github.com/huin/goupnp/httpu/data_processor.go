
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func processCSVFile(inputPath, outputPath string) error {
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

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    lineNumber := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV at line %d: %w", lineNumber, err)
        }

        cleanedRecord := make([]string, len(record))
        for i, field := range record {
            cleanedRecord[i] = strings.TrimSpace(field)
        }

        if len(cleanedRecord) > 0 && !isEmptyRecord(cleanedRecord) {
            if err := writer.Write(cleanedRecord); err != nil {
                return fmt.Errorf("error writing record at line %d: %w", lineNumber, err)
            }
        }

        lineNumber++
    }

    return nil
}

func isEmptyRecord(record []string) bool {
    for _, field := range record {
        if field != "" {
            return false
        }
    }
    return true
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("CSV file processed successfully")
}