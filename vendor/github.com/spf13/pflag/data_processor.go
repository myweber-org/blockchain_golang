
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
    Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        if len(row) != 4 {
            return nil, errors.New("invalid row format")
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(row[1])

        value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
        if err != nil {
            continue
        }

        valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

        record := DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
            Valid: valid,
        }
        records = append(records, record)
    }

    return records, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Valid && record.Value > 0 {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func CalculateAverage(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0
    }

    total := 0.0
    for _, record := range records {
        total += record.Value
    }
    return total / float64(len(records))
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        return
    }

    records, err := ParseCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error parsing file: %v\n", err)
        return
    }

    fmt.Printf("Total records: %d\n", len(records))

    validRecords := FilterValidRecords(records)
    fmt.Printf("Valid records: %d\n", len(validRecords))

    average := CalculateAverage(validRecords)
    fmt.Printf("Average value: %.2f\n", average)
}