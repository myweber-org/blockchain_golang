
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

func ValidateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = strings.Contains(record.Email, "@") && len(record.Name) > 0
		validated = append(validated, record)
	}
	return validated
}

func PrintRecords(records []DataRecord) {
	for _, record := range records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Status: %s\n",
			record.ID, record.Name, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob", "invalid-email", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	uniqueRecords := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(uniqueRecords)

	validatedRecords := ValidateRecords(uniqueRecords)
	fmt.Println("\nAfter validation:")
	PrintRecords(validatedRecords)
}