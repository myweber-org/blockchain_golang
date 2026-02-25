
package main

import (
	"encoding/csv"
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
	Active  bool
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
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

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
		}

		records = append(records, DataRecord{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		})
	}

	return records, nil
}

func validateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID %d found", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value %f for record ID %d", record.Value, record.ID)
		}
	}
	return nil
}

func processData(filename string) ([]DataRecord, error) {
	records, err := parseCSVFile(filename)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no data records found in file")
	}

	if err := validateRecords(records); err != nil {
		return nil, err
	}

	return records, nil
}

func calculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	for _, record := range records {
		sum += record.Value
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return sum, average, activeCount
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := processData(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		os.Exit(1)
	}

	total, average, activeCount := calculateStatistics(records)
	fmt.Printf("Processed %d records successfully\n", len(records))
	fmt.Printf("Total value: %.2f\n", total)
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Active records: %d\n", activeCount)
}