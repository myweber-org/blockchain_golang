package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID        int
	Name      string
	Value     float64
	Validated bool
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		if len(line) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		id, err := strconv.Atoi(strings.TrimSpace(line[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := strings.TrimSpace(line[1])
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(line[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		validated := strings.ToLower(strings.TrimSpace(line[3])) == "true"

		record := Record{
			ID:        id,
			Name:      name,
			Value:     value,
			Validated: validated,
		}

		records = append(records, record)
		lineNumber++
	}

	return records, nil
}

func calculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	validCount := 0

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
		if record.Validated {
			validCount++
		}
	}

	average := sum / float64(len(records))
	return average, max, validCount
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	average, max, validCount := calculateStatistics(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Maximum value: %.2f\n", max)
	fmt.Printf("Validated records: %d\n", validCount)
}