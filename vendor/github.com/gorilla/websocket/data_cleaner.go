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