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
package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}