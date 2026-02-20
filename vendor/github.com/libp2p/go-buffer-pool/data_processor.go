
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