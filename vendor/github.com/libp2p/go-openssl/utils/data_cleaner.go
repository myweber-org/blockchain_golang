
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	dc.Data = result
	return result
}

func (dc *DataCleaner) TrimWhitespace() []string {
	result := []string{}
	for _, item := range dc.Data {
		trimmed := strings.TrimSpace(item)
		result = append(result, trimmed)
	}
	dc.Data = result
	return result
}

func (dc *DataCleaner) Clean() []string {
	dc.TrimWhitespace()
	dc.RemoveDuplicates()
	return dc.Data
}

func main() {
	rawData := []string{"  apple ", "banana", "  apple", "cherry  ", "banana "}
	cleaner := NewDataCleaner(rawData)
	cleaned := cleaner.Clean()
	fmt.Println("Cleaned data:", cleaned)
}