
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; !exists {
			seen[trimmed] = struct{}{}
			result = append(result, trimmed)
		}
	}
	return result
}

func (dc *DataCleaner) ProcessLines(input []string) []string {
	return dc.RemoveDuplicates(input)
}

func main() {
	cleaner := &DataCleaner{}
	sample := []string{"  apple ", "banana", "apple", "  ", "banana", " cherry"}
	cleaned := cleaner.ProcessLines(sample)
	fmt.Println("Cleaned data:", cleaned)
}