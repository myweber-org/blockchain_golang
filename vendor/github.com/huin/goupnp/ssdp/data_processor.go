
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
    Email string
    Valid bool
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []DataRecord{}
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(line) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(line[0]),
            Name:  strings.TrimSpace(line[1]),
            Email: strings.TrimSpace(line[2]),
            Valid: validateRecord(strings.TrimSpace(line[0]), strings.TrimSpace(line[2])),
        }

        records = append(records, record)
    }

    return records, nil
}

func validateRecord(id, email string) bool {
    if id == "" || email == "" {
        return false
    }
    if !strings.Contains(email, "@") {
        return false
    }
    return true
}

func GenerateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid records: %d\n", validCount)
    fmt.Printf("Invalid records: %d\n", len(records)-validCount)
    
    if validCount > 0 {
        fmt.Println("\nValid records:")
        for _, record := range records {
            if record.Valid {
                fmt.Printf("  ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
            }
        }
    }
}