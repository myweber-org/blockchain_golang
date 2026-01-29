
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(row),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(row []string) bool {
	if len(row) < 3 {
		return false
	}

	id := strings.TrimSpace(row[0])
	name := strings.TrimSpace(row[1])
	email := strings.TrimSpace(row[2])

	if id == "" || name == "" || email == "" {
		return false
	}

	if !strings.Contains(email, "@") {
		return false
	}

	return true
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func GenerateReport(records []DataRecord) {
	validCount := 0
	for _, record := range records {
		if record.Valid {
			validCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", len(records)-validCount)

	if validCount > 0 {
		fmt.Println("\nValid Records:")
		for _, record := range records {
			if record.Valid {
				fmt.Printf("ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
			}
		}
	}
}