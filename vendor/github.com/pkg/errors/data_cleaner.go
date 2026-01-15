package utils

import "strings"

func SanitizeString(input string) string {
    trimmed := strings.TrimSpace(input)
    return strings.Join(strings.Fields(trimmed), " ")
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords int
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{processedRecords: 0}
}

func (dc *DataCleaner) RemoveDuplicates(data []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range data {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
			dc.processedRecords++
		}
	}
	return result
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) GetStats() string {
	return fmt.Sprintf("Processed %d unique records", dc.processedRecords)
}

func main() {
	cleaner := NewDataCleaner()
	
	sampleData := []string{
		"user@example.com",
		"  user@example.com  ",
		"invalid-email",
		"another@test.org",
		"",
		"another@test.org",
	}
	
	uniqueEmails := cleaner.RemoveDuplicates(sampleData)
	fmt.Println("Unique emails:", uniqueEmails)
	
	for _, email := range uniqueEmails {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("%s is valid\n", email)
		} else {
			fmt.Printf("%s is invalid\n", email)
		}
	}
	
	fmt.Println(cleaner.GetStats())
}