
package main

import (
	"fmt"
	"strings"
)

func CleanStringSlice(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	dirtyData := []string{"  apple ", "banana", "  apple", "banana ", "", "  cherry  "}
	cleaned := CleanStringSlice(dirtyData)
	fmt.Printf("Original: %v\n", dirtyData)
	fmt.Printf("Cleaned: %v\n", cleaned)
}