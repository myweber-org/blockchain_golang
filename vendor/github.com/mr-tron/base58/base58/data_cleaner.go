
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
		if !dc.seen[normalized] {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) Normalize(item string) string {
	return strings.ToUpper(strings.TrimSpace(item))
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "  Apple ", "banana", "BANANA", "cherry"}
	
	fmt.Println("Original:", data)
	
	unique := cleaner.RemoveDuplicates(data)
	fmt.Println("Unique:", unique)
	
	cleaner.Reset()
	
	for _, item := range data {
		fmt.Printf("Normalized '%s': %s\n", item, cleaner.Normalize(item))
	}
}