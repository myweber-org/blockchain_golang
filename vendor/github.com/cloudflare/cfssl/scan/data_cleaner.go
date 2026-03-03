package main

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
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
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
    cleaned := RemoveDuplicates(data)
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
    Phone string
}

type Cleaner struct {
    seenEmails map[string]bool
    seenPhones map[string]bool
}

func NewCleaner() *Cleaner {
    return &Cleaner{
        seenEmails: make(map[string]bool),
        seenPhones: make(map[string]bool),
    }
}

func (c *Cleaner) NormalizeEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func (c *Cleaner) NormalizePhone(phone string) string {
    return strings.ReplaceAll(strings.TrimSpace(phone), "-", "")
}

func (c *Cleaner) IsDuplicateEmail(email string) bool {
    normalized := c.NormalizeEmail(email)
    if c.seenEmails[normalized] {
        return true
    }
    c.seenEmails[normalized] = true
    return false
}

func (c *Cleaner) IsDuplicatePhone(phone string) bool {
    normalized := c.NormalizePhone(phone)
    if c.seenPhones[normalized] {
        return true
    }
    c.seenPhones[normalized] = true
    return false
}

func (c *Cleaner) ValidateEmail(email string) bool {
    normalized := c.NormalizeEmail(email)
    return strings.Contains(normalized, "@") && strings.Contains(normalized, ".")
}

func (c *Cleaner) ValidatePhone(phone string) bool {
    normalized := c.NormalizePhone(phone)
    return len(normalized) >= 10 && len(normalized) <= 15
}

func (c *Cleaner) ProcessRecords(records []DataRecord) []DataRecord {
    var cleaned []DataRecord
    for _, record := range records {
        if !c.ValidateEmail(record.Email) {
            continue
        }
        if !c.ValidatePhone(record.Phone) {
            continue
        }
        if c.IsDuplicateEmail(record.Email) {
            continue
        }
        if c.IsDuplicatePhone(record.Phone) {
            continue
        }
        cleaned = append(cleaned, record)
    }
    return cleaned
}

func main() {
    records := []DataRecord{
        {1, "test@example.com", "123-456-7890"},
        {2, "TEST@example.com", "1234567890"},
        {3, "invalid-email", "555-1234"},
        {4, "another@test.com", "123-456-7890"},
    }

    cleaner := NewCleaner()
    result := cleaner.ProcessRecords(records)

    fmt.Printf("Original records: %d\n", len(records))
    fmt.Printf("Cleaned records: %d\n", len(result))
    for _, r := range result {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", r.ID, r.Email, r.Phone)
    }
}