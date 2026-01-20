
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
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
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

func ProcessRecords(records []DataRecord, multiplier float64) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record, multiplier))
	}
	return processed, nil
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}

func FilterByTag(records []DataRecord, tag string) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		for _, t := range record.Tags {
			if t == tag {
				filtered = append(filtered, record)
				break
			}
		}
	}
	return filtered
}