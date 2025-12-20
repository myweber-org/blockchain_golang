
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
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []DataRecord{}
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

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        record, err := parseRecord(row, lineNumber)
        if err != nil {
            return nil, err
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func parseRecord(row []string, lineNumber int) (DataRecord, error) {
    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }

    name := strings.TrimSpace(row[1])
    if name == "" {
        return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
    }

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }

    return DataRecord{
        ID:    id,
        Name:  name,
        Value: value,
    }, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    min := records[0].Value
    max := records[0].Value

    for _, record := range records {
        sum += record.Value
        if record.Value < min {
            min = record.Value
        }
        if record.Value > max {
            max = record.Value
        }
    }

    average := sum / float64(len(records))
    return average, max - min, len(records)
}

func ValidateRecords(records []DataRecord) []error {
    var errors []error
    seenIDs := make(map[int]bool)

    for i, record := range records {
        if record.ID <= 0 {
            errors = append(errors, fmt.Errorf("record %d: invalid ID %d", i, record.ID))
        }

        if seenIDs[record.ID] {
            errors = append(errors, fmt.Errorf("record %d: duplicate ID %d", i, record.ID))
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            errors = append(errors, fmt.Errorf("record %d: negative value %f", i, record.Value))
        }
    }

    return errors
}