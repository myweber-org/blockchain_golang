package utils

func RemoveDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    
    return result
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	RemoveDuplicates bool
	TrimSpaces       bool
}

func NewDataCleaner(removeDuplicates, trimSpaces bool) *DataCleaner {
	return &DataCleaner{
		RemoveDuplicates: removeDuplicates,
		TrimSpaces:       trimSpaces,
	}
}

func (dc *DataCleaner) CleanStringSlice(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		processed := item
		if dc.TrimSpaces {
			processed = strings.TrimSpace(processed)
		}

		if processed == "" {
			continue
		}

		if dc.RemoveDuplicates {
			if seen[processed] {
				continue
			}
			seen[processed] = true
		}

		result = append(result, processed)
	}

	return result
}

func main() {
	data := []string{"  apple  ", "banana", "  apple  ", "cherry", "", "banana "}

	cleaner := NewDataCleaner(true, true)
	cleaned := cleaner.CleanStringSlice(data)

	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}