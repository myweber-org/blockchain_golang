package datautils

import "sort"

func RemoveDuplicates[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesSorted[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	sort.Slice(slice, func(i, j int) bool {
		// Convert to string for comparison to satisfy comparable constraint
		return fmt.Sprintf("%v", slice[i]) < fmt.Sprintf("%v", slice[j])
	})

	result := slice[:1]
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1] {
			result = append(result, slice[i])
		}
	}

	return result
}
package main

import (
	"fmt"
	"strings"
)

func deduplicateStrings(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func cleanData(input []string) []string {
	normalized := make([]string, len(input))
	for i, v := range input {
		normalized[i] = normalizeString(v)
	}
	return deduplicateStrings(normalized)
}

func main() {
	data := []string{"  Apple", "banana", "apple ", "Banana", "  Cherry  "}
	cleaned := cleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}