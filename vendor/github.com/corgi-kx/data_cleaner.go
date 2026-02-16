
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
	ID        int
	Name      string
	Email     string
	Age       int
	Active    bool
	Timestamp string
}

func parseCSVFile(filename string) ([]DataRecord, error) {
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

		if lineNumber == 1 {
			continue
		}

		if len(row) != 6 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 6, got %d", lineNumber, len(row))
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
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID: %w", err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("name cannot be empty")
	}

	record.Email = strings.TrimSpace(row[2])
	if !strings.Contains(record.Email, "@") {
		return record, fmt.Errorf("invalid email format")
	}

	record.Age, err = strconv.Atoi(strings.TrimSpace(row[3]))
	if err != nil || record.Age < 0 || record.Age > 150 {
		return record, fmt.Errorf("invalid age value")
	}

	record.Active, err = strconv.ParseBool(strings.TrimSpace(row[4]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag")
	}

	record.Timestamp = strings.TrimSpace(row[5])
	if record.Timestamp == "" {
		return record, fmt.Errorf("timestamp cannot be empty")
	}

	return record, nil
}

func validateRecords(records []DataRecord) ([]DataRecord, []string) {
	var validRecords []DataRecord
	var validationErrors []string

	for _, record := range records {
		if record.Age < 18 {
			validationErrors = append(validationErrors, fmt.Sprintf("Record ID %d: Age %d is below minimum requirement", record.ID, record.Age))
			continue
		}

		if !strings.HasSuffix(record.Email, ".com") && !strings.HasSuffix(record.Email, ".org") {
			validationErrors = append(validationErrors, fmt.Sprintf("Record ID %d: Email domain not supported", record.ID))
			continue
		}

		validRecords = append(validRecords, record)
	}

	return validRecords, validationErrors
}

func generateSummary(records []DataRecord) map[string]interface{} {
	if len(records) == 0 {
		return map[string]interface{}{
			"total_records":   0,
			"average_age":     0.0,
			"active_count":    0,
			"unique_domains":  0,
		}
	}

	totalAge := 0
	activeCount := 0
	domains := make(map[string]bool)

	for _, record := range records {
		totalAge += record.Age
		if record.Active {
			activeCount++
		}

		emailParts := strings.Split(record.Email, "@")
		if len(emailParts) == 2 {
			domains[emailParts[1]] = true
		}
	}

	return map[string]interface{}{
		"total_records":   len(records),
		"average_age":     float64(totalAge) / float64(len(records)),
		"active_count":    activeCount,
		"unique_domains":  len(domains),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run data_cleaner.go <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		os.Exit(1)
	}

	validRecords, errors := validateRecords(records)

	fmt.Printf("Processing completed:\n")
	fmt.Printf("  Total records read: %d\n", len(records))
	fmt.Printf("  Valid records: %d\n", len(validRecords))
	fmt.Printf("  Validation errors: %d\n", len(errors))

	if len(errors) > 0 {
		fmt.Println("\nValidation errors:")
		for _, errMsg := range errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	summary := generateSummary(validRecords)
	fmt.Println("\nData summary:")
	for key, value := range summary {
		fmt.Printf("  %s: %v\n", key, value)
	}
}