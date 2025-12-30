
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

    headerProcessed := false
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV record: %w", err)
        }

        if !headerProcessed {
            headerProcessed = true
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("error writing header: %w", err)
            }
            continue
        }

        cleanedRecord := cleanRecord(record)
        if isValidRecord(cleanedRecord) {
            if err := writer.Write(cleanedRecord); err != nil {
                return fmt.Errorf("error writing record: %w", err)
            }
        }
    }

    return nil
}

func cleanRecord(record []string) []string {
    cleaned := make([]string, len(record))
    for i, field := range record {
        cleaned[i] = strings.TrimSpace(field)
    }
    return cleaned
}

func isValidRecord(record []string) bool {
    for _, field := range record {
        if field == "" {
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

    inputPath := os.Args[1]
    outputPath := os.Args[2]

    if err := processCSVFile(inputPath, outputPath); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %s -> %s\n", inputPath, outputPath)
}