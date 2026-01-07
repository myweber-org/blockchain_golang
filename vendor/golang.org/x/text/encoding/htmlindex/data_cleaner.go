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

func (dc *DataCleaner) RemoveDuplicates(input []string) []string {
	var result []string
	for _, item := range input {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] {
			dc.seen[normalized] = true
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
	
	data := []string{
		"apple",
		"Apple",
		"banana",
		"  banana  ",
		"cherry",
		"APPLE",
		"date",
	}
	
	unique := cleaner.RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", unique)
	
	cleaner.Reset()
	testData := []string{"test", "TEST", "Test"}
	fmt.Println("After reset:", cleaner.RemoveDuplicates(testData))
}