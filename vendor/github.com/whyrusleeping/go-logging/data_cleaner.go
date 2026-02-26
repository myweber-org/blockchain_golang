
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Phone string
}

type Cleaner struct {
	seen map[string]bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		seen: make(map[string]bool),
	}
}

func (c *Cleaner) GenerateHash(r Record) string {
	data := fmt.Sprintf("%s-%s", strings.ToLower(r.Email), r.Phone)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (c *Cleaner) IsDuplicate(r Record) bool {
	hash := c.GenerateHash(r)
	if c.seen[hash] {
		return true
	}
	c.seen[hash] = true
	return false
}

func (c *Cleaner) ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (c *Cleaner) ValidatePhone(phone string) bool {
	if len(phone) < 10 {
		return false
	}
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func (c *Cleaner) ProcessRecords(records []Record) []Record {
	var cleaned []Record
	for _, rec := range records {
		if !c.ValidateEmail(rec.Email) || !c.ValidatePhone(rec.Phone) {
			continue
		}
		if c.IsDuplicate(rec) {
			continue
		}
		cleaned = append(cleaned, rec)
	}
	return cleaned
}

func main() {
	records := []Record{
		{"1", "test@example.com", "1234567890"},
		{"2", "test@example.com", "1234567890"},
		{"3", "invalid-email", "5555555555"},
		{"4", "valid@test.org", "9876543210"},
		{"5", "another@test.net", "1234567890"},
	}

	cleaner := NewCleaner()
	result := cleaner.ProcessRecords(records)

	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(result))
	for _, rec := range result {
		fmt.Printf("ID: %s, Email: %s\n", rec.ID, rec.Email)
	}
}