
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	line := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", line, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", line)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", line)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
		line++
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return fmt.Errorf("invalid ID %d", rec.ID)
		}
		if seen[rec.ID] {
			return fmt.Errorf("duplicate ID %d", rec.ID)
		}
		seen[rec.ID] = true
	}
	return nil
}

func CalculateTotal(records []Record) float64 {
	total := 0.0
	for _, rec := range records {
		total += rec.Value
	}
	return total
}