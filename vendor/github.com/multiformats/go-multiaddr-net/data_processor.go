
package main

import (
	"errors"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Tags  []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("ID must be positive")
	}
	if strings.TrimSpace(record.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Name = strings.ToUpper(strings.TrimSpace(record.Name))
	transformed.Value = record.Value * 1.1
	return transformed
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

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	var total float64
	for _, record := range records {
		total += record.Value
	}
	return total / float64(len(records))
}