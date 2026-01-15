package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Score int
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(id, email string, score int) {
	record := DataRecord{
		ID:    strings.TrimSpace(id),
		Email: strings.ToLower(strings.TrimSpace(email)),
		Score: score,
	}
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	unique := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := record.ID + "|" + record.Email
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}

	dc.records = unique
	return unique
}

func (dc *DataCleaner) ValidateRecords() (valid, invalid []DataRecord) {
	for _, record := range dc.records {
		if record.ID != "" && strings.Contains(record.Email, "@") && record.Score >= 0 && record.Score <= 100 {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func (dc *DataCleaner) GetAverageScore() float64 {
	if len(dc.records) == 0 {
		return 0.0
	}

	total := 0
	for _, record := range dc.records {
		total += record.Score
	}
	return float64(total) / float64(len(dc.records))
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord("001", "user1@example.com", 85)
	cleaner.AddRecord("002", "user2@example.com", 92)
	cleaner.AddRecord("001", "user1@example.com", 85)
	cleaner.AddRecord("003", "invalid-email", 105)
	cleaner.AddRecord("004", "user3@example.com", 78)

	fmt.Println("Original records:", len(cleaner.records))
	
	unique := cleaner.RemoveDuplicates()
	fmt.Println("After deduplication:", len(unique))

	valid, invalid := cleaner.ValidateRecords()
	fmt.Println("Valid records:", len(valid))
	fmt.Println("Invalid records:", len(invalid))

	average := cleaner.GetAverageScore()
	fmt.Printf("Average score: %.2f\n", average)
}
package main

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
}package datautils

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
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

func deduplicateEmails(emails []string) []string {
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

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailMap := make(map[string]bool)
	var validRecords []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailMap[cleanEmail] {
			emailMap[cleanEmail] = true
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "test@domain.org", false},
	}

	cleaned := processRecords(records)
	fmt.Printf("Processed %d records, %d valid after cleaning\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}