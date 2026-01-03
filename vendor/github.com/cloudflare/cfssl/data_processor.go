
package main

import (
    "encoding/csv"
    "errors"
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

func ParseCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    // Skip header
    _, err = reader.Read()
    if err != nil {
        return nil, err
    }

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) < 4 {
            continue
        }

        id, idErr := strconv.Atoi(strings.TrimSpace(row[0]))
        name := strings.TrimSpace(row[1])
        value, valErr := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
        valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

        if idErr != nil || valErr != nil || name == "" {
            continue
        }

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

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
    if len(records) == 0 {
        return nil, errors.New("no records to validate")
    }

    validRecords := make([]DataRecord, 0)
    for _, record := range records {
        if record.Valid && record.Value >= 0 {
            validRecords = append(validRecords, record)
        }
    }

    return validRecords, nil
}

func CalculateAverage(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0.0
    }

    total := 0.0
    for _, record := range records {
        total += record.Value
    }

    return total / float64(len(records))
}