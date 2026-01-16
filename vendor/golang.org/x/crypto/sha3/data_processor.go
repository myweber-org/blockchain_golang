
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append(record.Tags, "processed"),
	}
}

func ProcessData(records []DataRecord, multiplier float64) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record, multiplier))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Value:     42.5,
			Timestamp: time.Now(),
			Tags:      []string{"test"},
		},
		{
			ID:        "rec002",
			Value:     18.3,
			Timestamp: time.Now().Add(-time.Hour),
			Tags:      []string{"production"},
		},
	}

	processed, err := ProcessData(records, 2.0)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, rec := range processed {
		fmt.Printf("Processed: %+v\n", rec)
	}
}