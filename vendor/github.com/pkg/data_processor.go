package main

import (
    "errors"
    "fmt"
    "strings"
    "time"
)

type DataRecord struct {
    ID        string
    Timestamp time.Time
    Value     float64
    Tags      []string
}

func ValidateRecord(record DataRecord) error {
    if record.ID == "" {
        return errors.New("record ID cannot be empty")
    }
    if record.Value < 0 {
        return errors.New("record value cannot be negative")
    }
    if len(record.Tags) == 0 {
        return errors.New("record must have at least one tag")
    }
    return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
    return DataRecord{
        ID:        strings.ToUpper(record.ID),
        Timestamp: record.Timestamp.UTC(),
        Value:     record.Value * multiplier,
        Tags:      append([]string{"processed"}, record.Tags...),
    }
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
    var processed []DataRecord
    for _, record := range records {
        if err := ValidateRecord(record); err != nil {
            return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
        }
        processed = append(processed, TransformRecord(record, 1.5))
    }
    return processed, nil
}

func main() {
    records := []DataRecord{
        {ID: "rec001", Timestamp: time.Now(), Value: 100.0, Tags: []string{"test", "sample"}},
        {ID: "rec002", Timestamp: time.Now().Add(-time.Hour), Value: 200.0, Tags: []string{"production"}},
    }

    processed, err := ProcessRecords(records)
    if err != nil {
        fmt.Printf("Processing error: %v\n", err)
        return
    }

    fmt.Printf("Successfully processed %d records\n", len(processed))
    for _, rec := range processed {
        fmt.Printf("ID: %s, Value: %.2f, Tags: %v\n", rec.ID, rec.Value, rec.Tags)
    }
}