
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

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple", " BANANA ", "banana", "Cherry", "cherry "}
	
	fmt.Println("Original data:", data)
	
	uniqueData := cleaner.RemoveDuplicates(data)
	fmt.Println("Deduplicated data:", uniqueData)
	
	cleaner.Reset()
	
	testString := "  TEST Data  "
	normalized := cleaner.Normalize(testString)
	fmt.Printf("Normalized '%s' to '%s'\n", testString, normalized)
}