package main

import (
	"fmt"
	"strings"
)

func DeduplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func NormalizeWhitespace(input string) string {
	words := strings.Fields(input)
	return strings.Join(words, " ")
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana"}
	unique := DeduplicateStrings(data)
	fmt.Println("Deduplicated:", unique)

	text := "  Hello    world!  This   has  extra   spaces.  "
	normalized := NormalizeWhitespace(text)
	fmt.Println("Normalized:", normalized)
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

type DataCleaner struct {
    records []DataRecord
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        records: make([]DataRecord, 0),
    }
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
    dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
    seen := make(map[string]bool)
    unique := make([]DataRecord, 0)

    for _, record := range dc.records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }

    dc.records = unique
    return unique
}

func (dc *DataCleaner) ValidateEmails() []DataRecord {
    valid := make([]DataRecord, 0)

    for _, record := range dc.records {
        if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
            valid = append(valid, record)
        }
    }

    return valid
}

func (dc *DataCleaner) GetRecordCount() int {
    return len(dc.records)
}

func main() {
    cleaner := NewDataCleaner()

    sampleData := []DataRecord{
        {ID: 1, Email: "user@example.com", Phone: "1234567890"},
        {ID: 2, Email: "user@example.com", Phone: "1234567890"},
        {ID: 3, Email: "test@domain.org", Phone: "0987654321"},
        {ID: 4, Email: "invalid-email", Phone: "5555555555"},
        {ID: 5, Email: "another@test.net", Phone: "1111111111"},
    }

    for _, record := range sampleData {
        cleaner.AddRecord(record)
    }

    fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

    unique := cleaner.RemoveDuplicates()
    fmt.Printf("After deduplication: %d\n", len(unique))

    valid := cleaner.ValidateEmails()
    fmt.Printf("Valid email records: %d\n", len(valid))

    for _, record := range valid {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
    }
}