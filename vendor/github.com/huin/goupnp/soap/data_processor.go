
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

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
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

func NormalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var normalized []string
	
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		lower := strings.ToLower(trimmed)
		if trimmed != "" && !uniqueTags[lower] {
			uniqueTags[lower] = true
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
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