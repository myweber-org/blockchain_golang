package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes leading/trailing whitespace, reduces multiple spaces to single,
// and removes any non-printable characters from the input string.
func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned := spaceRegex.ReplaceAllString(trimmed, " ")
	
	// Remove non-printable characters
	printableRegex := regexp.MustCompile(`[^[:print:]]`)
	final := printableRegex.ReplaceAllString(cleaned, "")
	
	return final
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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
    return true
}

func CleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        record.Email = strings.ToLower(strings.TrimSpace(record.Email))
        if ValidateEmail(record.Email) && !emailSet[record.Email] {
            emailSet[record.Email] = true
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    testEmails := []string{
        "test@example.com",
        "TEST@example.com",
        "test@example.com",
        "invalid-email",
        "another@test.org",
    }
    
    fmt.Println("Original emails:", testEmails)
    deduped := DeduplicateEmails(testEmails)
    fmt.Println("Deduplicated emails:", deduped)
    
    records := []DataRecord{
        {1, "user1@domain.com", false},
        {2, "USER1@DOMAIN.COM", false},
        {3, "bad-email", false},
        {4, "user2@test.org", false},
    }
    
    cleaned := CleanData(records)
    fmt.Printf("Cleaned records: %+v\n", cleaned)
}