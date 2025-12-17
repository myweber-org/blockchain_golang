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

func RemoveDuplicates(records []DataRecord) []DataRecord {
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

func ValidateEmails(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if strings.Contains(email, "@") && strings.Contains(email, ".") {
			record.Valid = true
		} else {
			record.Valid = false
		}
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
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@domain.org", false},
		{5, "ANOTHER@DOMAIN.ORG", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	unique := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(unique)

	validated := ValidateEmails(unique)
	fmt.Println("\nAfter validation:")
	PrintRecords(validated)
}