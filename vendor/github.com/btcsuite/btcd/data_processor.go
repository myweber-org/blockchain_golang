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
package main

import "fmt"

func calculateMovingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || windowSize > len(data) {
        return []float64{}
    }

    result := make([]float64, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := i; j < i+windowSize; j++ {
            sum += data[j]
        }
        result[i] = sum / float64(windowSize)
    }
    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := calculateMovingAverage(sampleData, window)
    
    fmt.Printf("Data: %v\n", sampleData)
    fmt.Printf("Moving average (window=%d): %v\n", window, averages)
}