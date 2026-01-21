
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
	Valid     bool
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

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Tags = normalizeTags(record.Tags)
	transformed.Value = roundValue(record.Value)
	transformed.Valid = true
	return transformed
}

func normalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var result []string
	
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !uniqueTags[normalized] {
			uniqueTags[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func roundValue(value float64) float64 {
	return float64(int(value*100)) / 100
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		
		processed = append(processed, TransformRecord(record))
	}
	
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Timestamp: time.Now(),
			Value:     123.4567,
			Tags:      []string{"Sensor", "TEMPERATURE", "sensor"},
		},
		{
			ID:        "rec002",
			Timestamp: time.Now().Add(-1 * time.Hour),
			Value:     98.765,
			Tags:      []string{"pressure", "  PRESSURE  "},
		},
	}
	
	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	
	fmt.Printf("Successfully processed %d records\n", len(processed))
	for _, rec := range processed {
		fmt.Printf("Record %s: value=%.2f, tags=%v, valid=%v\n", 
			rec.ID, rec.Value, rec.Tags, rec.Valid)
	}
}