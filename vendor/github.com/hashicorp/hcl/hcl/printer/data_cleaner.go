
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	unique := DeduplicateRecords(records)

	for _, record := range unique {
		record.Valid = ValidateEmail(record.Email)
		if record.Valid {
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "ANOTHER@TEST.ORG", false},
	}

	cleaned := CleanData(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", r.ID, r.Email)
	}
}