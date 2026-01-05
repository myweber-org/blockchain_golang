
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(id int, email, name string) {
	record := DataRecord{
		ID:    id,
		Email: strings.ToLower(strings.TrimSpace(email)),
		Name:  strings.TrimSpace(name),
	}
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	unique := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%s|%s", record.Email, record.Name)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}

	dc.records = unique
	return unique
}

func (dc *DataCleaner) ValidateEmails() (valid, invalid []DataRecord) {
	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()
	
	cleaner.AddRecord(1, "user@example.com", "John Doe")
	cleaner.AddRecord(2, "user@example.com", "John Doe")
	cleaner.AddRecord(3, "jane.doe@domain.org", "Jane Doe")
	cleaner.AddRecord(4, "invalid-email", "Test User")
	cleaner.AddRecord(5, "another@test.com", "Another User")
	
	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())
	
	unique := cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", len(unique))
	
	valid, invalid := cleaner.ValidateEmails()
	fmt.Printf("Valid emails: %d, Invalid emails: %d\n", len(valid), len(invalid))
	
	for _, record := range valid {
		fmt.Printf("Valid: ID=%d, Email=%s, Name=%s\n", record.ID, record.Email, record.Name)
	}
}