
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Category  string
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
	if strings.TrimSpace(record.Category) == "" {
		return errors.New("category cannot be empty")
	}
	return nil
}

func TransformValue(value float64, multiplier float64) float64 {
	return value * multiplier
}

func NormalizeCategory(category string) string {
	return strings.ToUpper(strings.TrimSpace(category))
}

func ProcessRecord(record DataRecord, multiplier float64) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformedValue := TransformValue(record.Value, multiplier)
	normalizedCategory := NormalizeCategory(record.Category)

	return DataRecord{
		ID:        record.ID,
		Value:     transformedValue,
		Timestamp: record.Timestamp,
		Category:  normalizedCategory,
	}, nil
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}