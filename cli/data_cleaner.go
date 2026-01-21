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
		normalized := dc.normalize(item)
		if !dc.seen[normalized] {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple", " BANANA", "banana ", "Cherry", "cherry"}
	cleaned := cleaner.RemoveDuplicates(data)
	
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	
	cleaner.Reset()
	anotherSet := []string{"dog", "Dog", "CAT", "cat"}
	cleaned2 := cleaner.RemoveDuplicates(anotherSet)
	fmt.Println("Second set cleaned:", cleaned2)
}