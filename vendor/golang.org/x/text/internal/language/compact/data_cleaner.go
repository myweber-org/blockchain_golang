
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

func (dc *DataCleaner) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := strings.ToLower(trimmed)
	return normalized
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	cleaned := dc.CleanString(value)
	if dc.seen[cleaned] {
		return true
	}
	dc.seen[cleaned] = true
	return false
}

func (dc *DataCleaner) AddItem(value string) bool {
	if dc.IsDuplicate(value) {
		return false
	}
	return true
}

func (dc *DataCleaner) UniqueCount() int {
	return len(dc.seen)
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()

	samples := []string{"  Apple  ", "apple", "BANANA", "banana ", "  Cherry"}

	fmt.Println("Processing items:")
	for _, item := range samples {
		cleaned := cleaner.CleanString(item)
		added := cleaner.AddItem(item)
		fmt.Printf("Original: '%s' -> Cleaned: '%s' -> Added: %v\n", item, cleaned, added)
	}

	fmt.Printf("\nTotal unique items: %d\n", cleaner.UniqueCount())
}