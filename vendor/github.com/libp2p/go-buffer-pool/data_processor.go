
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

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 columns, got %d", lineNumber, len(line))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid ID: %v", lineNumber, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("line %d: name cannot be empty", lineNumber)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid value: %v", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNumber++
	}

	return records, nil
}

func calculateTotal(records []Record) float64 {
	total := 0.0
	for _, r := range records {
		total += r.Value
	}
	return total
}

func findMaxRecord(records []Record) *Record {
	if len(records) == 0 {
		return nil
	}
	max := records[0]
	for _, r := range records[1:] {
		if r.Value > max.Value {
			max = r
		}
	}
	return &max
}