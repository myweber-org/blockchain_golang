package datautils

func RemoveDuplicates(input []string) []string {
    seen := make(map[string]struct{})
    result := make([]string, 0, len(input))
    
    for _, item := range input {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
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

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && !strings.ContainsAny(item, "!@#$%")
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"apple",
		"Apple",
		"banana",
		"",
		"banana",
		"cherry!",
		"cherry",
		"  DATE  ",
		"date",
	}
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	
	cleaner.Reset()
	fmt.Println("Cleaner reset complete")
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

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	unique := DeduplicateRecords(records)

	for _, record := range unique {
		record.Valid = ValidateEmail(record.Email)
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "MIXED@EXAMPLE.COM", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Processed %d records\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Phone string
	Tags  []string
}

type Cleaner struct {
	seenIDs    map[string]bool
	validEmails map[string]bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		seenIDs:    make(map[string]bool),
		validEmails: make(map[string]bool),
	}
}

func (c *Cleaner) GenerateID(email, phone string) string {
	combined := fmt.Sprintf("%s|%s", strings.ToLower(email), phone)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:8])
}

func (c *Cleaner) IsDuplicate(record DataRecord) bool {
	if record.ID == "" {
		record.ID = c.GenerateID(record.Email, record.Phone)
	}
	
	if c.seenIDs[record.ID] {
		return true
	}
	
	c.seenIDs[record.ID] = true
	return false
}

func (c *Cleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	
	if c.validEmails[email] {
		return false
	}
	
	c.validEmails[email] = true
	return true
}

func (c *Cleaner) NormalizeTags(tags []string) []string {
	unique := make(map[string]bool)
	var result []string
	
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !unique[normalized] {
			unique[normalized] = true
			result = append(result, normalized)
		}
	}
	
	return result
}

func (c *Cleaner) ProcessRecords(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	
	for _, record := range records {
		if c.IsDuplicate(record) {
			continue
		}
		
		if !c.ValidateEmail(record.Email) {
			continue
		}
		
		record.Tags = c.NormalizeTags(record.Tags)
		cleaned = append(cleaned, record)
	}
	
	return cleaned
}

func main() {
	cleaner := NewCleaner()
	
	records := []DataRecord{
		{Email: "test@example.com", Phone: "1234567890", Tags: []string{"Go", "go", "  "}},
		{Email: "test@example.com", Phone: "1234567890", Tags: []string{"Python"}},
		{Email: "invalid-email", Phone: "0987654321", Tags: []string{"Java"}},
		{Email: "new@domain.com", Phone: "5551234567", Tags: []string{"Rust", "rust", "Systems"}},
	}
	
	cleaned := cleaner.ProcessRecords(records)
	
	fmt.Printf("Original: %d records\n", len(records))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	
	for i, record := range cleaned {
		fmt.Printf("Record %d: %s - Tags: %v\n", i+1, record.Email, record.Tags)
	}
}