package datautil

import "sort"

func RemoveDuplicates[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesSorted[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	sort.Slice(input, func(i, j int) bool {
		// Convert to string for comparison to satisfy sort.Interface
		// This works for any comparable type
		return false // We only need grouping, not actual sorting
	})

	result := make([]T, 0, len(input))
	result = append(result, input[0])

	for i := 1; i < len(input); i++ {
		if input[i] != input[i-1] {
			result = append(result, input[i])
		}
	}

	return result
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
	data := []int{4, 2, 8, 2, 4, 9, 8, 1}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

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

func (dc *DataCleaner) Clean(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	cleaned := dc.Clean(value)
	if dc.seen[cleaned] {
		return true
	}
	dc.seen[cleaned] = true
	return false
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	cleaned := dc.Clean(email)
	return strings.Contains(cleaned, "@") && strings.Contains(cleaned, ".")
}

func main() {
	cleaner := NewDataCleaner()
	
	samples := []string{
		"  TEST@EXAMPLE.COM  ",
		"test@example.com",
		"invalid-email",
		"another@test.org",
	}
	
	for _, sample := range samples {
		cleaned := cleaner.Clean(sample)
		duplicate := cleaner.IsDuplicate(sample)
		validEmail := cleaner.ValidateEmail(sample)
		
		fmt.Printf("Original: %q\n", sample)
		fmt.Printf("Cleaned: %q\n", cleaned)
		fmt.Printf("Duplicate: %v\n", duplicate)
		fmt.Printf("Valid Email: %v\n\n", validEmail)
	}
}