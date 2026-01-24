package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)
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

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID format at line %d: %w", lineNumber, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value format at line %d: %w", lineNumber, err)
        }

        if value < 0 {
            return nil, errors.New("negative values are not allowed")
        }

        records = append(records, DataRecord{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid data records found in file")
    }

    return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, record := range records {
        sum += record.Value
    }
    average := sum / float64(len(records))

    var varianceSum float64
    for _, record := range records {
        diff := record.Value - average
        varianceSum += diff * diff
    }
    variance := varianceSum / float64(len(records))

    return average, variance
}