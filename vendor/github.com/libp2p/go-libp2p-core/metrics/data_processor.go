package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessRecord(record DataRecord) (string, error) {
	if record.ID <= 0 {
		return "", errors.New("invalid record ID")
	}

	if strings.TrimSpace(record.Name) == "" {
		return "", errors.New("record name cannot be empty")
	}

	if !record.Valid {
		return "", errors.New("record is marked as invalid")
	}

	processedValue := record.Value * 1.1
	result := fmt.Sprintf("Processed record %d: %s -> %.2f", record.ID, record.Name, processedValue)

	return result, nil
}

func ValidateAndProcess(records []DataRecord) ([]string, []error) {
	var results []string
	var errs []error

	for _, record := range records {
		result, err := ProcessRecord(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", record.ID, err))
			continue
		}
		results = append(results, result)
	}

	return results, errs
}

func main() {
	records := []DataRecord{
		{ID: 1, Name: "Record One", Value: 100.0, Valid: true},
		{ID: 2, Name: "", Value: 200.0, Valid: true},
		{ID: 0, Name: "Record Three", Value: 300.0, Valid: true},
		{ID: 4, Name: "Record Four", Value: 400.0, Valid: false},
		{ID: 5, Name: "Record Five", Value: 500.0, Valid: true},
	}

	results, errs := ValidateAndProcess(records)

	fmt.Println("Processing Results:")
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("\nErrors:")
	for _, err := range errs {
		fmt.Println(err)
	}
}