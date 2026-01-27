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

func ValidateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = strings.Contains(records[i].Email, "@") &&
			len(records[i].Email) > 3 &&
			strings.Contains(records[i].Email, ".")
	}
	return records
}

func CleanData(records []DataRecord) []DataRecord {
	records = RemoveDuplicates(records)
	records = ValidateEmails(records)
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob", "invalid-email", false},
	}

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Valid: %v\n", record.ID, record.Name, record.Valid)
	}
}