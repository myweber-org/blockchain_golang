
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

func (dc *DataCleaner) Deduplicate(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && !strings.ContainsAny(item, "!@#$%")
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "Apple", "banana", "", "cherry!", "banana", "date"}
	cleaned := cleaner.Deduplicate(data)
	
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "GRAPE", "kiwi"}
	moreCleaned := cleaner.Deduplicate(moreData)
	fmt.Println("More cleaned:", moreCleaned)
}
package main

import "fmt"

func removeDuplicates(input []int) []int {
    seen := make(map[int]bool)
    result := []int{}

    for _, value := range input {
        if !seen[value] {
            seen[value] = true
            result = append(result, value)
        }
    }
    return result
}

func main() {
    slice := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := removeDuplicates(slice)
    fmt.Println("Original:", slice)
    fmt.Println("Cleaned:", cleaned)
}