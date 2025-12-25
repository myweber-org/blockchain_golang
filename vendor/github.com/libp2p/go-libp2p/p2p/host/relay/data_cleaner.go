
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

func deduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    var unique []DataRecord
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validatePhone(phone string) bool {
    if len(phone) != 10 {
        return false
    }
    for _, ch := range phone {
        if ch < '0' || ch > '9' {
            return false
        }
    }
    return true
}

func cleanData(records []DataRecord) []DataRecord {
    var cleaned []DataRecord
    uniqueRecords := deduplicateRecords(records)
    
    for _, record := range uniqueRecords {
        if validateEmail(record.Email) && validatePhone(record.Phone) {
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    sampleData := []DataRecord{
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 2, Email: "invalid-email", Phone: "9876543210"},
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 3, Email: "user@domain.org", Phone: "5551234"},
    }
    
    cleaned := cleanData(sampleData)
    fmt.Printf("Original records: %d\n", len(sampleData))
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
    }
}