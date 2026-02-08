
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) < 3 {
			return nil, errors.New("invalid row format: insufficient columns")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		}

		if err := validateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %d: %w", id, err)
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("ID must be positive")
	}
	if record.Name == "" {
		return errors.New("name cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	return nil
}

func CalculateTotalValue(records []DataRecord) float64 {
	total := 0.0
	for _, record := range records {
		total += record.Value
	}
	return total
}

func FindMaxValueRecord(records []DataRecord) *DataRecord {
	if len(records) == 0 {
		return nil
	}

	maxRecord := records[0]
	for _, record := range records[1:] {
		if record.Value > maxRecord.Value {
			maxRecord = record
		}
	}
	return &maxRecord
}
package main

import (
    "errors"
    "strings"
)

func ValidateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}

func TrimAndTitle(s string) string {
    trimmed := strings.TrimSpace(s)
    return strings.Title(strings.ToLower(trimmed))
}

func FilterEmptyStrings(slice []string) []string {
    result := []string{}
    for _, s := range slice {
        if s != "" {
            result = append(result, s)
        }
    }
    return result
}