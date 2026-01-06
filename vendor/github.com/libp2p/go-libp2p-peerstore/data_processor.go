
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
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:    strings.TrimSpace(row[0]),
            Name:  strings.TrimSpace(row[1]),
            Value: strings.TrimSpace(row[2]),
        }

        if record.ID == "" || record.Name == "" {
            return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) error {
    idSet := make(map[string]bool)
    
    for _, record := range records {
        if idSet[record.ID] {
            return fmt.Errorf("duplicate ID found: %s", record.ID)
        }
        idSet[record.ID] = true
        
        if len(record.Value) > 100 {
            return fmt.Errorf("value too long for record ID %s", record.ID)
        }
    }
    
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    records, err := ProcessCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    if err := ValidateRecords(records); err != nil {
        fmt.Printf("Validation error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records\n", len(records))
    for _, record := range records {
        fmt.Printf("ID: %s, Name: %s, Value: %s\n", record.ID, record.Name, record.Value)
    }
}