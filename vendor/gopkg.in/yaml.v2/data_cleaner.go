
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

func DeduplicateEmails(emails []string) []string {
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

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        if ValidateEmail(record.Email) && !emailSet[record.Email] {
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
    
    cleaned := CleanData(records)
    fmt.Printf("Original: %d records\n", len(records))
    fmt.Printf("Cleaned: %d records\n", len(cleaned))
    
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
}