package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, rec := range records {
		email := strings.ToLower(strings.TrimSpace(rec.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, rec)
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func CleanData(records []Record) ([]Record, error) {
	var cleaned []Record
	for _, rec := range records {
		if err := ValidateEmail(rec.Email); err != nil {
			continue
		}
		cleaned = append(cleaned, rec)
	}
	cleaned = DeduplicateRecords(cleaned)
	if len(cleaned) == 0 {
		return cleaned, errors.New("no valid records after cleaning")
	}
	return cleaned, nil
}

func main() {
	sampleData := []Record{
		{1, "user@example.com", true},
		{2, "invalid-email", true},
		{3, "user@example.com", true},
		{4, "another@test.org", true},
		{5, "", true},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Cleaned %d records\n", len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", rec.ID, rec.Email)
	}
}