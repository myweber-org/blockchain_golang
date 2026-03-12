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
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

		records = append(records, Record{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		})
	}

	return records, nil
}

func validateRecords(records []Record) error {
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

func processData(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	totalValue := 0.0
	activeCount := 0
	for _, record := range records {
		totalValue += record.Value
		if record.Active {
			activeCount++
		}
	}

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", totalValue)
	fmt.Printf("Active records: %d\n", activeCount)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	if err := processData(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Value float64
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(record.Email) {
		return errors.New("invalid email format")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	return DataRecord{
		ID:    strings.ToUpper(record.ID),
		Email: strings.ToLower(record.Email),
		Value: record.Value * 1.1,
	}
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, err
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}
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
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", lineNum)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return records, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return fmt.Errorf("invalid ID %d", rec.ID)
		}
		if rec.Name == "" {
			return fmt.Errorf("empty name for ID %d", rec.ID)
		}
		if rec.Value < 0 {
			return fmt.Errorf("negative value for ID %d", rec.ID)
		}
		if seenIDs[rec.ID] {
			return fmt.Errorf("duplicate ID %d", rec.ID)
		}
		seenIDs[rec.ID] = true
	}
	return nil
}

func calculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	min := records[0].Value
	max := records[0].Value

	for _, rec := range records {
		sum += rec.Value
		if rec.Value < min {
			min = rec.Value
		}
		if rec.Value > max {
			max = rec.Value
		}
	}

	average := sum / float64(len(records))
	return average, max - min
}