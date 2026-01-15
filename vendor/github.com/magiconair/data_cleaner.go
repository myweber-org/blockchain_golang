
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

func (dc *DataCleaner) Clean(input []string) []string {
	var result []string
	for _, item := range input {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && normalized != "" {
			dc.seen[normalized] = true
			result = append(result, normalized)
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
		"  Apple  ",
		"apple",
		"BANANA",
		"  banana  ",
		"",
		"Cherry",
		"cherry ",
	}
	
	cleaned := cleaner.Clean(data)
	fmt.Println("Cleaned data:", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "GRAPE", "kiwi"}
	moreCleaned := cleaner.Clean(moreData)
	fmt.Println("More cleaned data:", moreCleaned)
}