
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

func (dc *DataCleaner) AddToSet(value string) {
	cleaned := dc.CleanString(value)
	dc.seen[cleaned] = true
}

func (dc *DataCleaner) GetUniqueCount() int {
	return len(dc.seen)
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	samples := []string{"  Apple  ", "apple", "BANANA", "banana ", "  Cherry"}
	
	for _, sample := range samples {
		if cleaner.IsDuplicate(sample) {
			fmt.Printf("Duplicate found: %s\n", sample)
		} else {
			fmt.Printf("New unique value: %s\n", sample)
			cleaner.AddToSet(sample)
		}
	}
	
	fmt.Printf("Total unique values: %d\n", cleaner.GetUniqueCount())
}