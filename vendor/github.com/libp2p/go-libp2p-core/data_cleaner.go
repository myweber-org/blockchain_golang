
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

func DeduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] && len(email) > 0 {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	seenEmails := make(map[string]bool)

	for _, record := range records {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		if ValidateEmail(record.Email) && !seenEmails[record.Email] {
			record.Valid = true
			seenEmails[record.Email] = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "  USER@EXAMPLE.COM  ", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "", false},
	}

	cleaned := CleanRecords(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}