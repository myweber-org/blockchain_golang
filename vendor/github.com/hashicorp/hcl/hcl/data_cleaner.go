
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

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validateRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		record.Valid = validateEmail(record.Email)
		valid = append(valid, record)
	}
	return valid
}

func processData(records []DataRecord) []DataRecord {
	deduped := deduplicateRecords(records)
	validated := validateRecords(deduped)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "user@example.com", false},
	}

	processed := processData(sampleData)

	for _, record := range processed {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
	}
}
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
	seen := make(map[string]struct{})
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if _, exists := seen[email]; !exists {
			seen[email] = struct{}{}
			result = append(result, email)
		}
	}
	return result
}

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func CleanRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if ValidateEmail(cleanEmail) && !emailSet[cleanEmail] {
			emailSet[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.com", false},
		{5, "test@domain.com", false},
	}

	cleaned := CleanRecords(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}

	emails := []string{"a@b.com", "A@B.COM", "  a@b.com  ", "invalid"}
	unique := DeduplicateEmails(emails)
	fmt.Println("Unique emails:", unique)
}