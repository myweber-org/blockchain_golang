
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
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		email := records[i].Email
		records[i].Valid = strings.Contains(email, "@") &&
			strings.Contains(email, ".") &&
			len(email) > 5
	}
	return records
}

func PrintRecords(records []DataRecord) {
	for _, record := range records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n",
			record.ID, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "   User@Example.Com   ", false},
		{6, "", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	unique := RemoveDuplicates(records)
	validated := ValidateEmails(unique)

	fmt.Println("\nCleaned and validated records:")
	PrintRecords(validated)
}