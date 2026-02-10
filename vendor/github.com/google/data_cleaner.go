
package main

import (
	"fmt"
	"sort"
)

type DataRecord struct {
	ID   int
	Name string
}

func CleanData(records []DataRecord) []DataRecord {
	seen := make(map[int]bool)
	var unique []DataRecord

	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].ID < unique[j].ID
	})

	return unique
}

func main() {
	records := []DataRecord{
		{ID: 3, Name: "Charlie"},
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 1, Name: "AliceDuplicate"},
		{ID: 4, Name: "David"},
	}

	cleaned := CleanData(records)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
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

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func ProcessRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if ValidateEmail(record.Email) {
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return DeduplicateRecords(validRecords)
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "test@domain.org", false},
		{5, "another@test.co.uk", false},
	}

	processed := ProcessRecords(records)
	fmt.Printf("Processed %d records\n", len(processed))
	for _, r := range processed {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}