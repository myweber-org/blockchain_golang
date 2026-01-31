package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        int
	Value     string
	Timestamp time.Time
	Valid     bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid record ID")
	}
	if strings.TrimSpace(record.Value) == "" {
		return errors.New("record value cannot be empty")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformValue(input string) string {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) > 50 {
		trimmed = trimmed[:50] + "..."
	}
	return strings.ToUpper(trimmed)
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %d: %w", record.ID, err)
		}
		record.Value = TransformValue(record.Value)
		record.Valid = true
		processed = append(processed, record)
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{ID: 1, Value: "sample data", Timestamp: time.Now().Add(-time.Hour)},
		{ID: 2, Value: "   ", Timestamp: time.Now().Add(-time.Hour)},
	}
	
	result, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	
	fmt.Printf("Processed %d records successfully\n", len(result))
}