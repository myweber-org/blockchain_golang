
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
		record.Valid = record.ID > 0 &&
			len(strings.TrimSpace(record.Name)) > 0 &&
			strings.Contains(record.Email, "@")
		validated = append(validated, record)
	}
	return validated
}

func PrintRecords(records []DataRecord) {
	fmt.Println("ID\tName\tEmail\t\tValid")
	fmt.Println("----------------------------------------")
	for _, record := range records {
		fmt.Printf("%d\t%s\t%s\t%t\n", record.ID, record.Name, record.Email, record.Valid)
	}
}

func main() {
	records := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{0, "", "invalid-email", false},
		{4, "Alice Brown", "alice@example.com", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	uniqueRecords := RemoveDuplicates(records)
	fmt.Println("\nAfter removing duplicates:")
	PrintRecords(uniqueRecords)

	validatedRecords := ValidateRecords(uniqueRecords)
	fmt.Println("\nAfter validation:")
	PrintRecords(validatedRecords)
}