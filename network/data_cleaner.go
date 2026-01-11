
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

func processRecords(records []DataRecord) []DataRecord {
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
        {4, "test@domain.org", false},
    }

    cleaned := processRecords(records)
    fmt.Printf("Processed %d records\n", len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
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

func (dc *DataCleaner) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ReplaceAll(trimmed, "\"", "'")
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
	
	for _, record := range deduped {
		fmt.Printf("Email '%s' valid: %v\n", 
			record, cleaner.ValidateEmail(record))
	}
	
	testInput := `  This has "quotes" and spaces  `
	fmt.Printf("Sanitized: '%s'\n", cleaner.SanitizeInput(testInput))
}