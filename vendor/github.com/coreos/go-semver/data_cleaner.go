
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

func (dc *DataCleaner) CleanPhoneNumber(phone string) string {
	var builder strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func main() {
	cleaner := NewDataCleaner()

	// Test deduplication
	records := []string{"John Doe", "john doe", "Jane Smith", "JOHN DOE"}
	unique := cleaner.RemoveDuplicates(records)
	fmt.Printf("Unique records: %v\n", unique)

	// Test email validation
	emails := []string{"test@example.com", "invalid", "user@domain"}
	for _, email := range emails {
		fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
	}

	// Test phone cleaning
	phone := "+1 (555) 123-4567"
	cleaned := cleaner.CleanPhoneNumber(phone)
	fmt.Printf("Cleaned phone: %s\n", cleaned)
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
    numbers := []int{1, 2, 2, 3, 4, 4, 5, 5, 5}
    uniqueNumbers := RemoveDuplicates(numbers)
    fmt.Println("Original:", numbers)
    fmt.Println("Unique:", uniqueNumbers)
    
    strings := []string{"apple", "banana", "apple", "orange", "banana"}
    uniqueStrings := RemoveDuplicates(strings)
    fmt.Println("Original:", strings)
    fmt.Println("Unique:", uniqueStrings)
}