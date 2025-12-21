
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]Record, 0)

    for line := 1; ; line++ {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", line, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d", line)
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found")
    }

    return records, nil
}

func CalculateStats(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, r := range records {
        sum += r.Value
    }
    average := sum / float64(len(records))

    var variance float64
    for _, r := range records {
        diff := r.Value - average
        variance += diff * diff
    }
    stdDev := variance / float64(len(records))

    return average, stdDev
}

func ValidateRecord(r Record) error {
    if r.ID <= 0 {
        return errors.New("ID must be positive")
    }
    if r.Name == "" {
        return errors.New("name cannot be empty")
    }
    if r.Value < 0 {
        return errors.New("value cannot be negative")
    }
    return nil
}