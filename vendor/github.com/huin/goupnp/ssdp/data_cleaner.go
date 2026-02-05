package main

import "fmt"

func removeDuplicates(input []int) []int {
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    var result []DataRecord
    for _, record := range records {
        normalizedEmail := strings.ToLower(strings.TrimSpace(record.Email))
        if !seen[normalizedEmail] {
            seen[normalizedEmail] = true
            result = append(result, record)
        }
    }
    return result
}

func ValidateEmails(records []DataRecord) []DataRecord {
    for i := range records {
        email := records[i].Email
        records[i].Valid = strings.Contains(email, "@") && strings.Contains(email, ".")
    }
    return records
}

func ProcessData(input []DataRecord) []DataRecord {
    deduped := RemoveDuplicates(input)
    validated := ValidateEmails(deduped)
    return validated
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "user@example.com", false},
        {3, "invalid-email", false},
        {4, "test@domain.org", false},
    }

    processed := ProcessData(sampleData)
    fmt.Printf("Processed %d records\n", len(processed))
    for _, record := range processed {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
    }
}