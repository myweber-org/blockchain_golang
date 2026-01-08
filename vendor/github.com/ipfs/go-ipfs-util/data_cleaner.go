package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID      int
    Name    string
    Email   string
    Active  bool
    Score   float64
}

func CleanCSVData(filename string) ([]DataRecord, error) {
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

        if len(row) < 5 {
            continue
        }

        record, err := parseRow(row)
        if err != nil {
            fmt.Printf("skipping invalid row at line %d: %v\n", lineNumber, err)
            continue
        }

        if validateRecord(record) {
            records = append(records, record)
        }
    }

    return records, nil
}

func parseRow(row []string) (DataRecord, error) {
    var record DataRecord
    var err error

    record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID: %w", err)
    }

    record.Name = strings.TrimSpace(row[1])
    if record.Name == "" {
        return record, fmt.Errorf("name cannot be empty")
    }

    record.Email = strings.TrimSpace(row[2])
    if !strings.Contains(record.Email, "@") {
        return record, fmt.Errorf("invalid email format")
    }

    record.Active, err = strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, fmt.Errorf("invalid active status: %w", err)
    }

    record.Score, err = strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid score: %w", err)
    }

    return record, nil
}

func validateRecord(record DataRecord) bool {
    if record.ID <= 0 {
        return false
    }
    if record.Score < 0 || record.Score > 100 {
        return false
    }
    if len(record.Email) > 254 {
        return false
    }
    return true
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_cleaner <csv_file>")
        os.Exit(1)
    }

    records, err := CleanCSVData(os.Args[1])
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d valid records\n", len(records))
    for _, r := range records {
        fmt.Printf("ID: %d, Name: %s, Email: %s, Active: %v, Score: %.2f\n",
            r.ID, r.Name, r.Email, r.Active, r.Score)
    }
}