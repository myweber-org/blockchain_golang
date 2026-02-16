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
	seen := make(map[string]struct{})
	result := []string{}
	for _, email := range emails {
		if _, exists := seen[email]; !exists {
			seen[email] = struct{}{}
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord

	for _, record := range records {
		if validateEmail(record.Email) && !emailSet[record.Email] {
			emailSet[record.Email] = true
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "another@test.org", false},
	}

	cleaned := cleanData(records)
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}
package main

import (
	"fmt"
	"strings"
)

func deduplicateStrings(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func cleanData(input []string) []string {
	normalized := make([]string, len(input))
	for i, v := range input {
		normalized[i] = normalizeString(v)
	}
	return deduplicateStrings(normalized)
}

func main() {
	sample := []string{"  Apple", "banana", "apple", " Banana ", "cherry"}
	cleaned := cleanData(sample)
	fmt.Println("Original:", sample)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{
		"user@example.com",
		"USER@example.com",
		"test@domain.org",
		"invalid-email",
		"test@domain.org",
	}
	
	fmt.Println("Original records:", records)
	unique := cleaner.RemoveDuplicates(records)
	fmt.Println("Deduplicated records:", unique)
	
	for _, record := range unique {
		if cleaner.ValidateEmail(record) {
			fmt.Printf("Valid email: %s\n", record)
		} else {
			fmt.Printf("Invalid email: %s\n", record)
		}
	}
}