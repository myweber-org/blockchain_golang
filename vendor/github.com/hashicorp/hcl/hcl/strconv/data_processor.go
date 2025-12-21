
package main

import (
	"encoding/csv"
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil {
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row length: %d", len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
		}

		record := DataRecord{
			ID:    id,
			Name:  row[1],
			Value: value,
		}

		if !validateRecord(record) {
			return nil, fmt.Errorf("validation failed for record ID: %d", id)
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID <= 0 {
		return false
	}
	if record.Name == "" {
		return false
	}
	if record.Value < 0 {
		return false
	}
	return true
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}