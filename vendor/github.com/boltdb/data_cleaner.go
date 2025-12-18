
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		record.Valid = strings.Contains(record.Email, "@") && len(record.Name) > 0
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func processData(records []DataRecord) []DataRecord {
	records = deduplicateRecords(records)
	records = validateRecords(records)
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Invalid User", "invalid-email", false},
	}

	cleanedData := processData(sampleData)
	fmt.Printf("Processed %d records\n", len(cleanedData))
	for _, record := range cleanedData {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}