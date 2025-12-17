
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

func (dc *DataCleaner) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	replacer := strings.NewReplacer("\n", " ", "\t", " ", "\r", " ")
	return replacer.Replace(trimmed)
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{
		"user@example.com",
		"  USER@EXAMPLE.COM  ",
		"invalid-email",
		"another@test.org",
		"user@example.com",
	}
	
	fmt.Println("Original records:", records)
	deduped := cleaner.RemoveDuplicates(records)
	fmt.Println("After deduplication:", deduped)
	
	for _, email := range deduped {
		sanitized := cleaner.SanitizeInput(email)
		valid := cleaner.ValidateEmail(sanitized)
		fmt.Printf("Email: %s, Valid: %v\n", sanitized, valid)
	}
}