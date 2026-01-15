
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	result := []DataRecord{}

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}
	return result
}

func ValidateRecords(records []DataRecord) []DataRecord {
	validated := []DataRecord{}
	for _, record := range records {
		record.Valid = record.ID > 0 &&
			len(strings.TrimSpace(record.Name)) > 0 &&
			strings.Contains(record.Email, "@")
		validated = append(validated, record)
	}
	return validated
}

func CleanData(records []DataRecord) []DataRecord {
	deduped := RemoveDuplicates(records)
	validated := ValidateRecords(deduped)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{0, "Invalid User", "invalid-email", false},
	}

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Valid: %v\n",
			record.ID, record.Name, record.Email, record.Valid)
	}
}