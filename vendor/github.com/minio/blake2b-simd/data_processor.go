
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID    string
	Name  string
	Email string
	Valid bool
}

func processCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(line[0]),
			Name:  strings.TrimSpace(line[1]),
			Email: strings.TrimSpace(line[2]),
			Valid: validateRecord(strings.TrimSpace(line[0]), strings.TrimSpace(line[2])),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func generateReport(records []DataRecord) {
	validCount := 0
	invalidCount := 0

	fmt.Println("=== DATA PROCESSING REPORT ===")
	for _, record := range records {
		if record.Valid {
			validCount++
			fmt.Printf("✓ Valid: %s - %s\n", record.ID, record.Name)
		} else {
			invalidCount++
			fmt.Printf("✗ Invalid: %s - %s (Email: %s)\n", record.ID, record.Name, record.Email)
		}
	}

	fmt.Printf("\nSummary: %d total records\n", len(records))
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", invalidCount)
	fmt.Printf("Validation rate: %.1f%%\n", float64(validCount)/float64(len(records))*100)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run data_processor.go <csv_file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	records, err := processCSVFile(filePath)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	generateReport(records)
}