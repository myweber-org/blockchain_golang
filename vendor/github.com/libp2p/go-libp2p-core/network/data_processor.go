
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

type Record struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

type DataProcessor struct {
    records []Record
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        records: make([]Record, 0),
    }
}

func (dp *DataProcessor) LoadFromCSV(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    lineNumber := 0
    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("csv read error on line %d: %w", lineNumber, err)
        }

        if len(row) != 4 {
            return fmt.Errorf("invalid column count on line %d: expected 4, got %d", lineNumber, len(row))
        }

        record, err := parseRecord(row)
        if err != nil {
            return fmt.Errorf("parse error on line %d: %w", lineNumber, err)
        }

        dp.records = append(dp.records, record)
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var record Record

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID format: %w", err)
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, errors.New("name cannot be empty")
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid value format: %w", err)
    }
    record.Value = value

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, fmt.Errorf("invalid active flag format: %w", err)
    }
    record.Active = active

    return record, nil
}

func (dp *DataProcessor) FilterActive() []Record {
    var activeRecords []Record
    for _, record := range dp.records {
        if record.Active {
            activeRecords = append(activeRecords, record)
        }
    }
    return activeRecords
}

func (dp *DataProcessor) CalculateTotal() float64 {
    var total float64
    for _, record := range dp.records {
        total += record.Value
    }
    return total
}

func (dp *DataProcessor) FindByName(name string) *Record {
    for _, record := range dp.records {
        if strings.EqualFold(record.Name, name) {
            return &record
        }
    }
    return nil
}

func (dp *DataProcessor) ExportToCSV(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, record := range dp.records {
        row := []string{
            strconv.Itoa(record.ID),
            record.Name,
            strconv.FormatFloat(record.Value, 'f', 2, 64),
            strconv.FormatBool(record.Active),
        }
        if err := writer.Write(row); err != nil {
            return fmt.Errorf("write error: %w", err)
        }
    }

    return nil
}

func (dp *DataProcessor) PrintSummary() {
    fmt.Printf("Total records: %d\n", len(dp.records))
    fmt.Printf("Active records: %d\n", len(dp.FilterActive()))
    fmt.Printf("Total value: %.2f\n", dp.CalculateTotal())
}