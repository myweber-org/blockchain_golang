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

func DeduplicateRecords(records []DataRecord) []DataRecord {
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

func ValidateRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		record.Valid = strings.Contains(record.Email, "@") && len(record.Name) > 0
		if record.Valid {
			valid = append(valid, record)
		}
	}
	return valid
}

func CleanDataPipeline(records []DataRecord) []DataRecord {
	unique := DeduplicateRecords(records)
	valid := ValidateRecords(unique)
	return valid
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Invalid User", "invalid-email", false},
	}

	cleaned := CleanDataPipeline(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", r.ID, r.Name, r.Email)
	}
}