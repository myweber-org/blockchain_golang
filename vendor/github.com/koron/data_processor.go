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

type Record struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if lineNum == 0 {
			lineNum++
			continue
		}

		record, err := parseLine(line, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
		lineNum++
	}

	return records, nil
}

func parseLine(fields []string, lineNum int) (Record, error) {
	if len(fields) != 4 {
		return Record{}, fmt.Errorf("invalid field count at line %d: expected 4, got %d", lineNum, len(fields))
	}

	id, err := strconv.Atoi(strings.TrimSpace(fields[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}

	name := strings.TrimSpace(fields[1])
	if name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}

	valid := strings.ToLower(strings.TrimSpace(fields[3])) == "true"

	return Record{
		ID:    id,
		Name:  name,
		Value: value,
		Valid: valid,
	}, nil
}

func ValidateRecords(records []Record) ([]Record, []error) {
	validRecords := []Record{}
	errorsList := []error{}

	seenIDs := make(map[int]bool)

	for _, record := range records {
		var recordErrors []string

		if record.ID <= 0 {
			recordErrors = append(recordErrors, "ID must be positive")
		}
		if seenIDs[record.ID] {
			recordErrors = append(recordErrors, "duplicate ID")
		}
		if len(record.Name) > 100 {
			recordErrors = append(recordErrors, "name exceeds 100 characters")
		}
		if record.Value < 0 {
			recordErrors = append(recordErrors, "value cannot be negative")
		}

		if len(recordErrors) > 0 {
			errorsList = append(errorsList, fmt.Errorf("record ID %d errors: %s", 
				record.ID, strings.Join(recordErrors, ", ")))
		} else {
			validRecords = append(validRecords, record)
			seenIDs[record.ID] = true
		}
	}

	return validRecords, errorsList
}

func CalculateStatistics(records []Record) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records to process")
	}

	var sum float64
	var validCount int

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			validCount++
		}
	}

	if validCount == 0 {
		return 0, 0, errors.New("no valid records found")
	}

	average := sum / float64(validCount)

	var varianceSum float64
	for _, record := range records {
		if record.Valid {
			diff := record.Value - average
			varianceSum += diff * diff
		}
	}
	variance := varianceSum / float64(validCount)

	return average, variance, nil
}