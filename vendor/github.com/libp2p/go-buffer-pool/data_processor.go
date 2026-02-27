
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
package data

import (
	"regexp"
	"strings"
)

// SanitizeInput removes extra whitespace and validates string content
func SanitizeInput(input string) (string, error) {
	// Trim leading/trailing whitespace
	cleaned := strings.TrimSpace(input)
	
	// Check for empty string after trimming
	if cleaned == "" {
		return "", ErrEmptyInput
	}
	
	// Remove multiple consecutive spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")
	
	// Validate against basic injection patterns (simplified example)
	injectionPattern := regexp.MustCompile(`(?i)(select|insert|delete|update|drop|union)`)
	if injectionPattern.MatchString(cleaned) {
		return "", ErrInvalidInput
	}
	
	return cleaned, nil
}

// ValidateEmail checks if a string is a valid email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Errors
var (
	ErrEmptyInput   = errors.New("input cannot be empty")
	ErrInvalidInput = errors.New("input contains invalid characters")
)
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
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNum, len(row))
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
			return fmt.Errorf("invalid record ID: %d must be positive", rec.ID)
		}
		if rec.Name == "" {
			return fmt.Errorf("record ID %d has empty name", rec.ID)
		}
		if rec.Value < 0 {
			return fmt.Errorf("record ID %d has negative value: %f", rec.ID, rec.Value)
		}
		if seenIDs[rec.ID] {
			return fmt.Errorf("duplicate record ID: %d", rec.ID)
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
	var max float64 = records[0].Value

	for _, rec := range records {
		sum += rec.Value
		if rec.Value > max {
			max = rec.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}