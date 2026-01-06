
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
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if strings.Contains(record.Email, "@") && len(record.Email) > 3 {
			record.Valid = true
			valid = append(valid, record)
		}
	}
	return valid
}

func cleanDataPipeline(records []DataRecord) []DataRecord {
	deduped := deduplicateRecords(records)
	validated := validateEmails(deduped)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "admin@test.org", false},
		{4, "invalid-email", false},
		{5, "ADMIN@TEST.ORG", false},
	}

	cleaned := cleanDataPipeline(sampleData)

	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}
package main

import "fmt"

func removeDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	cleanedData := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleanedData)
}