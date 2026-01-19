
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

func ValidateEmail(email string) bool {
    if len(email) < 3 || !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    return strings.Contains(parts[1], ".")
}

func CleanData(records []DataRecord) []DataRecord {
    emailMap := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
        if ValidateEmail(cleanEmail) && !emailMap[cleanEmail] {
            emailMap[cleanEmail] = true
            record.Email = cleanEmail
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    emails := []string{
        "test@example.com",
        "TEST@example.com",
        "invalid-email",
        "another@test.org",
        "test@example.com",
        "  spaced@email.com  ",
    }
    
    fmt.Println("Original emails:", emails)
    fmt.Println("Deduplicated:", DeduplicateEmails(emails))
    
    records := []DataRecord{
        {1, "user@domain.com", false},
        {2, "DUPLICATE@domain.com", false},
        {3, "duplicate@domain.com", false},
        {4, "bad-email", false},
        {5, "new@test.co.uk", false},
    }
    
    cleaned := CleanData(records)
    fmt.Printf("\nCleaned records: %d out of %d\n", len(cleaned), len(records))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
}