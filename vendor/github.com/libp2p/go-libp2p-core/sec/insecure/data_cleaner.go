
package main

import (
    "fmt"
    "strings"
)

func DeduplicateStrings(slice []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, item := range slice {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
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
    if !strings.Contains(parts[1], ".") {
        return false
    }
    return true
}

func CleanPhoneNumber(phone string) string {
    var builder strings.Builder
    for _, ch := range phone {
        if ch >= '0' && ch <= '9' {
            builder.WriteRune(ch)
        }
    }
    return builder.String()
}

func main() {
    emails := []string{"test@example.com", "invalid-email", "test@example.com", "another@domain.org"}
    fmt.Println("Original emails:", emails)
    uniqueEmails := DeduplicateStrings(emails)
    fmt.Println("Deduplicated emails:", uniqueEmails)

    for _, email := range uniqueEmails {
        fmt.Printf("Email %s is valid: %v\n", email, ValidateEmail(email))
    }

    phone := "+1 (234) 567-8900"
    cleaned := CleanPhoneNumber(phone)
    fmt.Printf("Original phone: %s -> Cleaned: %s\n", phone, cleaned)
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
			record.Valid = true
			valid = append(valid, record)
		}
	}
	return valid
}

func ProcessData(records []DataRecord) []DataRecord {
	deduped := DeduplicateRecords(records)
	validated := ValidateEmails(deduped)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "ANOTHER@TEST.ORG", false},
	}

	processed := ProcessData(sampleData)
	fmt.Printf("Original: %d, Processed: %d\n", len(sampleData), len(processed))
	for _, record := range processed {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
}