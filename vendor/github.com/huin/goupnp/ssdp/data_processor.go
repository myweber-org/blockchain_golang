
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID    string
    Name  string
    Email string
    Valid bool
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []DataRecord{}
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(line) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(line[0]),
            Name:  strings.TrimSpace(line[1]),
            Email: strings.TrimSpace(line[2]),
            Valid: validateRecord(strings.TrimSpace(line[0]), strings.TrimSpace(line[2])),
        }

        records = append(records, record)
    }

    return records, nil
}

func validateRecord(id, email string) bool {
    if id == "" || email == "" {
        return false
    }
    if !strings.Contains(email, "@") {
        return false
    }
    return true
}

func GenerateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid records: %d\n", validCount)
    fmt.Printf("Invalid records: %d\n", len(records)-validCount)
    
    if validCount > 0 {
        fmt.Println("\nValid records:")
        for _, record := range records {
            if record.Valid {
                fmt.Printf("  ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
            }
        }
    }
}
package data

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
	ErrEmptyData    = errors.New("data cannot be empty")
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return ErrInvalidInput
	}
	if record.Value < 0 {
		return ErrInvalidInput
	}
	if record.Timestamp.IsZero() {
		return ErrInvalidInput
	}
	return nil
}

func TransformTags(tags []string) []string {
	var transformed []string
	for _, tag := range tags {
		cleanTag := strings.TrimSpace(tag)
		if cleanTag != "" {
			transformed = append(transformed, strings.ToLower(cleanTag))
		}
	}
	return transformed
}

func CalculateAverage(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptyData
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values)), nil
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

func GenerateReport(records []DataRecord) map[string]interface{} {
	if len(records) == 0 {
		return map[string]interface{}{
			"count":   0,
			"message": "no records to process",
		}
	}

	var totalValue float64
	tagSet := make(map[string]bool)

	for _, record := range records {
		totalValue += record.Value
		for _, tag := range record.Tags {
			tagSet[tag] = true
		}
	}

	uniqueTags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		uniqueTags = append(uniqueTags, tag)
	}

	return map[string]interface{}{
		"count":        len(records),
		"total_value":  totalValue,
		"unique_tags":  len(uniqueTags),
		"average":      totalValue / float64(len(records)),
		"sample_tags":  uniqueTags[:min(3, len(uniqueTags))],
		"processed_at": time.Now().UTC(),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}