
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
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	unique := DeduplicateRecords(records)

	for _, record := range unique {
		record.Valid = ValidateEmail(record.Email)
		if record.Valid {
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "ANOTHER@TEST.ORG", false},
	}

	cleaned := CleanData(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", r.ID, r.Email)
	}
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID   int
	Name string
	Age  int
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		key := fmt.Sprintf("%d-%s-%d", record.ID, strings.ToLower(record.Name), record.Age)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateRecord(record DataRecord) error {
	if record.ID <= 0 {
		return fmt.Errorf("invalid ID: %d", record.ID)
	}
	if strings.TrimSpace(record.Name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if record.Age < 0 || record.Age > 150 {
		return fmt.Errorf("age out of valid range: %d", record.Age)
	}
	return nil
}

func cleanData(records []DataRecord) ([]DataRecord, []string) {
	var cleaned []DataRecord
	var errors []string

	uniqueRecords := deduplicateRecords(records)

	for _, record := range uniqueRecords {
		if err := validateRecord(record); err != nil {
			errors = append(errors, fmt.Sprintf("Record ID %d: %v", record.ID, err))
			continue
		}
		cleaned = append(cleaned, record)
	}

	return cleaned, errors
}

func main() {
	records := []DataRecord{
		{1, "John Doe", 30},
		{2, "Jane Smith", 25},
		{1, "John Doe", 30},
		{3, "", 40},
		{4, "Bob Johnson", -5},
		{5, "Alice Brown", 200},
	}

	cleaned, errors := cleanData(records)

	fmt.Println("Cleaned Records:")
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", record.ID, record.Name, record.Age)
	}

	if len(errors) > 0 {
		fmt.Println("\nValidation Errors:")
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}