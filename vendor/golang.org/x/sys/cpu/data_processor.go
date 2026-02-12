
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
	IsValid bool
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
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

		if len(row) < 3 {
			continue
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			fmt.Printf("Warning line %d: %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID format: %s", row[0])
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("name cannot be empty")
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %s", row[2])
	}
	record.Value = value

	record.IsValid = record.Value > 0 && record.ID > 0

	return record, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var total float64
	validCount := 0

	for _, record := range records {
		if record.IsValid {
			total += record.Value
			validCount++
		}
	}

	if validCount == 0 {
		return 0, 0, 0
	}

	average := total / float64(validCount)
	return total, average, validCount
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file_path>")
		return
	}

	filePath := os.Args[1]
	records, err := ProcessCSVFile(filePath)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	total, average, validCount := CalculateStatistics(records)

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Total value: %.2f\n", total)
	fmt.Printf("Average value: %.2f\n", average)
}