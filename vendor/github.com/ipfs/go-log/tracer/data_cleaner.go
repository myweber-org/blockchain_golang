
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple", " BANANA ", "banana", "Cherry", "cherry "}
	
	fmt.Println("Original data:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated data:", deduped)
	
	cleaner.Reset()
	testValue := "TestValue"
	fmt.Printf("Is '%s' duplicate? %v\n", testValue, cleaner.IsDuplicate(testValue))
	fmt.Printf("Is '%s' duplicate? %v\n", strings.ToLower(testValue), cleaner.IsDuplicate(strings.ToLower(testValue)))
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
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
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
	emailMap := make(map[string]bool)
	var validRecords []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailMap[cleanEmail] {
			emailMap[cleanEmail] = true
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "user@example.com", false},
	}

	cleaned := processRecords(records)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}