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

func ProcessCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]Record, 0)

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        if len(row) != 3 {
            return nil, errors.New("invalid row format")
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID: %w", err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value: %w", err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    return records, nil
}

func ValidateRecords(records []Record) error {
    if len(records) == 0 {
        return errors.New("no records to validate")
    }

    seen := make(map[int]bool)
    for _, r := range records {
        if r.ID <= 0 {
            return fmt.Errorf("invalid ID %d", r.ID)
        }
        if r.Name == "" {
            return fmt.Errorf("empty name for ID %d", r.ID)
        }
        if r.Value < 0 {
            return fmt.Errorf("negative value for ID %d", r.ID)
        }
        if seen[r.ID] {
            return fmt.Errorf("duplicate ID %d", r.ID)
        }
        seen[r.ID] = true
    }

    return nil
}

func CalculateTotalValue(records []Record) float64 {
    total := 0.0
    for _, r := range records {
        total += r.Value
    }
    return total
}