
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

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord
	
	for _, rec := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(rec.Email))
		if validateEmail(cleanEmail) && !emailSet[cleanEmail] {
			emailSet[cleanEmail] = true
			rec.Email = cleanEmail
			rec.Valid = true
			cleaned = append(cleaned, rec)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "test@domain.org", false},
	}
	
	cleaned := processRecords(records)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(records), len(cleaned))
	
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", rec.ID, rec.Email, rec.Valid)
	}
}