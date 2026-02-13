
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	dc.Data = result
	return result
}

func (dc *DataCleaner) TrimWhitespace() []string {
	result := []string{}
	for _, item := range dc.Data {
		trimmed := strings.TrimSpace(item)
		result = append(result, trimmed)
	}
	dc.Data = result
	return result
}

func (dc *DataCleaner) Clean() []string {
	dc.TrimWhitespace()
	dc.RemoveDuplicates()
	return dc.Data
}

func main() {
	rawData := []string{"  apple ", "banana", "  apple", "cherry  ", "banana "}
	cleaner := NewDataCleaner(rawData)
	cleaned := cleaner.Clean()
	fmt.Println("Cleaned data:", cleaned)
}
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

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
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
	for _, record := range records {
		if err := ValidateEmail(record.Email); err != nil {
			continue
		}
		cleaned = append(cleaned, record)
	}
	return DeduplicateRecords(cleaned), nil
}

func main() {
	records := []Record{
		{1, "user@example.com", true},
		{2, "USER@example.com", true},
		{3, "test@domain.org", true},
		{4, "invalid-email", true},
		{5, "", true},
	}

	cleaned, err := CleanData(records)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original: %d records\n", len(records))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", r.ID, r.Email)
	}
}