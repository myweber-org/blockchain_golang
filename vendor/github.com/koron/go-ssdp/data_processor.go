
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	average, max, count := CalculateStatistics(records)
	fmt.Printf("Processed %d records\n", count)
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Maximum value: %.2f\n", max)
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	userData := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	userData = TransformUsername(userData)

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}

func main() {
	processedData, err := ProcessUserInput("  JohnDoe  ", "john@example.com", 30)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}
package data_processor

import (
	"regexp"
	"strings"
)

type Processor struct {
	whitespaceRegex *regexp.Regexp
}

func NewProcessor() *Processor {
	return &Processor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (p *Processor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := p.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (p *Processor) NormalizeCase(input string) string {
	return strings.ToLower(input)
}

func (p *Processor) ExtractTokens(input string) []string {
	cleaned := p.CleanInput(input)
	normalized := p.NormalizeCase(cleaned)
	return strings.Fields(normalized)
}package main

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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func ValidateRecords(records []DataRecord) []error {
	var errors []error

	for i, record := range records {
		if record.ID <= 0 {
			errors = append(errors, fmt.Errorf("record %d: ID must be positive", i+1))
		}
		if len(record.Name) == 0 {
			errors = append(errors, fmt.Errorf("record %d: name cannot be empty", i+1))
		}
		if record.Value < 0 {
			errors = append(errors, fmt.Errorf("record %d: value cannot be negative", i+1))
		}
	}

	return errors
}

func ProcessData(filename string) error {
	records, err := ParseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	validationErrors := ValidateRecords(records)
	if len(validationErrors) > 0 {
		fmt.Println("Validation errors found:")
		for _, err := range validationErrors {
			fmt.Printf("  - %v\n", err)
		}
		return fmt.Errorf("data validation failed")
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	for _, record := range records {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f, Active: %t\n",
			record.ID, record.Name, record.Value, record.Active)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	if err := ProcessData(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}