
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
		if !seen[email] && isValidEmail(email) {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validateRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.ID > 0 && isValidEmail(record.Email) {
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	emails := []string{
		"user@example.com",
		"USER@example.com",
		"test@domain.org",
		"invalid-email",
		"user@example.com",
	}

	uniqueEmails := deduplicateEmails(emails)
	fmt.Println("Deduplicated emails:", uniqueEmails)

	records := []DataRecord{
		{ID: 1, Email: "alice@example.com"},
		{ID: 2, Email: "bob@test.org"},
		{ID: 0, Email: "invalid@com"},
		{ID: 3, Email: "charlie@domain.net"},
	}

	validRecords := validateRecords(records)
	fmt.Println("Valid records count:", len(validRecords))
	for _, r := range validRecords {
		fmt.Printf("ID: %d, Email: %s\n", r.ID, r.Email)
	}
}