
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    result := []DataRecord{}
    
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            result = append(result, record)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanPhoneNumber(phone string) string {
    cleaned := strings.Builder{}
    for _, ch := range phone {
        if ch >= '0' && ch <= '9' {
            cleaned.WriteRune(ch)
        }
    }
    return cleaned.String()
}

func ProcessRecords(records []DataRecord) []DataRecord {
    uniqueRecords := RemoveDuplicates(records)
    
    for i := range uniqueRecords {
        uniqueRecords[i].Phone = CleanPhoneNumber(uniqueRecords[i].Phone)
    }
    
    validRecords := []DataRecord{}
    for _, record := range uniqueRecords {
        if ValidateEmail(record.Email) {
            validRecords = append(validRecords, record)
        }
    }
    
    return validRecords
}

func main() {
    sampleData := []DataRecord{
        {ID: 1, Email: "test@example.com", Phone: "(123) 456-7890"},
        {ID: 2, Email: "invalid-email", Phone: "555-1234"},
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 3, Email: "user@domain.org", Phone: "+1-800-555-0199"},
    }
    
    cleaned := ProcessRecords(sampleData)
    
    fmt.Printf("Processed %d records\n", len(cleaned))
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", 
            record.ID, record.Email, record.Phone)
    }
}