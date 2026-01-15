
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
	uniqueRecords := []string{}
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}
	return uniqueRecords
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
	cleaned := strings.Builder{}
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			cleaned.WriteRune(ch)
		}
	}
	return cleaned.String()
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{
		"user@example.com",
		"  USER@EXAMPLE.COM  ",
		"invalid-email",
		"user@example.com",
		"another@test.org",
	}
	
	fmt.Println("Original records:", records)
	deduped := cleaner.RemoveDuplicates(records)
	fmt.Println("After deduplication:", deduped)
	
	testEmails := []string{"test@domain.com", "bad-email", "a@b.c"}
	for _, email := range testEmails {
		fmt.Printf("Email '%s' valid: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	phone := "+1 (234) 567-8900"
	fmt.Printf("Cleaned phone '%s': %s\n", phone, cleaner.CleanPhoneNumber(phone))
}