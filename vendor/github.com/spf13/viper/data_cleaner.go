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

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, DataRecord{
				ID:    record.ID,
				Email: email,
				Valid: record.Valid,
			})
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		isValid := strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".")
		validated = append(validated, DataRecord{
			ID:    record.ID,
			Email: record.Email,
			Valid: isValid,
		})
	}
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "  user@example.com  ", false},
	}

	fmt.Println("Original records:", len(sampleData))
	deduped := DeduplicateRecords(sampleData)
	fmt.Println("After deduplication:", len(deduped))

	validated := ValidateEmails(deduped)
	validCount := 0
	for _, r := range validated {
		if r.Valid {
			validCount++
		}
	}
	fmt.Println("Valid emails:", validCount)
}package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}