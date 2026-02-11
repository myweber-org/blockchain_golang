
package main

import (
	"fmt"
	"sort"
)

type DataRecord struct {
	ID   int
	Name string
}

type DataSet []DataRecord

func (d DataSet) RemoveDuplicates() DataSet {
	seen := make(map[int]bool)
	result := DataSet{}
	for _, record := range d {
		if !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

func (d DataSet) SortByID() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})
}

func CleanData(records DataSet) DataSet {
	unique := records.RemoveDuplicates()
	unique.SortByID()
	return unique
}

func main() {
	data := DataSet{
		{ID: 5, Name: "ItemE"},
		{ID: 2, Name: "ItemB"},
		{ID: 5, Name: "ItemE"},
		{ID: 1, Name: "ItemA"},
		{ID: 2, Name: "ItemB"},
	}

	cleaned := CleanData(data)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
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
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = strings.Contains(records[i].Email, "@") &&
			len(records[i].Email) > 3
	}
	return records
}

func CleanDataPipeline(records []DataRecord) []DataRecord {
	records = DeduplicateRecords(records)
	records = ValidateEmails(records)
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
	}

	cleaned := CleanDataPipeline(sampleData)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", r.ID, r.Email, r.Valid)
	}
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateEmails(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, Record{
				ID:    len(unique) + 1,
				Email: email,
				Valid: record.Valid,
			})
		}
	}
	return unique
}

func ValidateEmailFormat(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @ symbol")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errors.New("invalid email format")
	}
	if !strings.Contains(parts[1], ".") {
		return errors.New("domain must contain a dot")
	}
	return nil
}

func CleanRecords(records []Record) ([]Record, error) {
	cleaned := DeduplicateEmails(records)
	for i := range cleaned {
		if err := ValidateEmailFormat(cleaned[i].Email); err != nil {
			cleaned[i].Valid = false
		} else {
			cleaned[i].Valid = true
		}
	}
	return cleaned, nil
}

func main() {
	sampleData := []Record{
		{1, "user@example.com", true},
		{2, "USER@example.com", true},
		{3, "invalid-email", true},
		{4, "another@test.org", true},
		{5, "", true},
	}

	cleaned, err := CleanRecords(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Cleaned Records:")
	for _, record := range cleaned {
		status := "Valid"
		if !record.Valid {
			status = "Invalid"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
	}
}