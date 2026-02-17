
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

func (dc *DataCleaner) Clean(input string) string {
	normalized := strings.ToLower(strings.TrimSpace(input))
	if dc.seen[normalized] {
		return ""
	}
	dc.seen[normalized] = true
	return normalized
}

func (dc *DataCleaner) ProcessBatch(items []string) []string {
	var result []string
	for _, item := range items {
		cleaned := dc.Clean(item)
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"  Apple ", "apple", "BANANA", "banana ", "  Cherry  "}
	cleaned := cleaner.ProcessBatch(data)
	fmt.Println("Cleaned data:", cleaned)
}