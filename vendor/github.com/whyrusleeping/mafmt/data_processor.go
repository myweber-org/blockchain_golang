package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID    string
    Name  string
    Value string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []DataRecord
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(row[0]),
            Name:  strings.TrimSpace(row[1]),
            Value: strings.TrimSpace(row[2]),
        }

        if record.ID == "" || record.Name == "" {
            continue
        }

        records = append(records, record)
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    seenIDs := make(map[string]bool)

    for _, record := range records {
        if seenIDs[record.ID] {
            continue
        }
        if len(record.Value) > 0 && record.Value != "null" {
            seenIDs[record.ID] = true
            validRecords = append(validRecords, record)
        }
    }

    return validRecords
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        return
    }

    records, err := ProcessCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    validRecords := ValidateRecords(records)
    fmt.Printf("Processed %d records, %d valid records found\n", len(records), len(validRecords))
}