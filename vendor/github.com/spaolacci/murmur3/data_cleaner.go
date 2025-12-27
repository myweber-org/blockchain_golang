
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func (dc *DataCleaner) ProcessRecords(records []string) []string {
	cleaned := dc.RemoveDuplicates(records)
	return cleaned
}

func main() {
	cleaner := &DataCleaner{}
	data := []string{"apple", " banana", "apple", "cherry ", "", "banana", "  "}
	result := cleaner.ProcessRecords(data)
	fmt.Println("Cleaned data:", result)
}