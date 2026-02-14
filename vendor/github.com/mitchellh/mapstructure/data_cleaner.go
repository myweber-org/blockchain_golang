
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc DataCleaner) TrimWhitespace(items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func main() {
	cleaner := DataCleaner{}
	data := []string{" apple ", "banana", " apple ", "  cherry  ", "banana"}

	fmt.Println("Original data:", data)
	trimmed := cleaner.TrimWhitespace(data)
	fmt.Println("After trimming:", trimmed)
	deduped := cleaner.RemoveDuplicates(trimmed)
	fmt.Println("After deduplication:", deduped)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
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
	slice := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Cleaned:", cleaned)
}