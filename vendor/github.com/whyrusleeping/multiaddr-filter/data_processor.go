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
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record, err := parseRow(row)
		if err != nil {
			return nil, fmt.Errorf("parse error at line %d: %w", lineNumber, err)
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID format: %w", err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %w", err)
	}
	record.Value = value

	valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid valid flag format: %w", err)
	}
	record.Valid = valid

	return record, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []DataRecord) {
	var validRecords []DataRecord
	var invalidRecords []DataRecord

	for _, record := range records {
		if record.Valid && record.Value >= 0 && record.Name != "" {
			validRecords = append(validRecords, record)
		} else {
			invalidRecords = append(invalidRecords, record)
		}
	}

	return validRecords, invalidRecords
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 {
			min = record.Value
			max = record.Value
		} else {
			if record.Value < min {
				min = record.Value
			}
			if record.Value > max {
				max = record.Value
			}
		}
	}

	average := sum / float64(count)
	return average, min, max
}

func ProcessDataFile(filename string) error {
	records, err := ParseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	validRecords, invalidRecords := ValidateRecords(records)

	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Invalid records: %d\n", len(invalidRecords))

	if len(validRecords) > 0 {
		avg, min, max := CalculateStatistics(validRecords)
		fmt.Printf("Statistics for valid records:\n")
		fmt.Printf("  Average value: %.2f\n", avg)
		fmt.Printf("  Minimum value: %.2f\n", min)
		fmt.Printf("  Maximum value: %.2f\n", max)
	}

	return nil
}