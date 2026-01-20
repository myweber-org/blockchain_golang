
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	IsValid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("no valid data records found")
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNumber)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	isValid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid boolean at line %d: %w", lineNumber, err)
	}
	record.IsValid = isValid

	return record, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var validCount int
	var maxValue float64

	for _, record := range records {
		if record.IsValid {
			sum += record.Value
			validCount++
			if record.Value > maxValue {
				maxValue = record.Value
			}
		}
	}

	average := 0.0
	if validCount > 0 {
		average = sum / float64(validCount)
	}

	return average, maxValue, validCount
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.IsValid {
			filtered = append(filtered, record)
		}
	}
	return filtered
}