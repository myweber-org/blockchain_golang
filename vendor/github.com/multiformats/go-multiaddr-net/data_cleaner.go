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
		if !dc.seen[normalized] && len(normalized) > 0 {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"test@example.com", "  TEST@example.com  ", "invalid", "another@test.org", ""}
	
	fmt.Println("Original data:", data)
	
	uniqueData := cleaner.RemoveDuplicates(data)
	fmt.Println("After deduplication:", uniqueData)
	
	for _, item := range uniqueData {
		if cleaner.ValidateEmail(item) {
			fmt.Printf("Valid email: %s\n", item)
		} else {
			fmt.Printf("Invalid email: %s\n", item)
		}
	}
}