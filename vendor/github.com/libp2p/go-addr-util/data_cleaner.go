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
		normalized := strings.ToLower(strings.TrimSpace(email))
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func validateEmailFormat(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	cleaned := []DataRecord{}
	for _, rec := range records {
		normalizedEmail := strings.ToLower(strings.TrimSpace(rec.Email))
		if !emailSet[normalizedEmail] && validateEmailFormat(normalizedEmail) {
			emailSet[normalizedEmail] = true
			rec.Email = normalizedEmail
			rec.Valid = true
			cleaned = append(cleaned, rec)
		}
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "another@domain.org", false},
		{5, "user@example.com", false},
	}
	cleaned := cleanData(sampleData)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(sampleData), len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", rec.ID, rec.Email, rec.Valid)
	}
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}