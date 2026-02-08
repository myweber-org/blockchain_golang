package main

import (
	"fmt"
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func CleanData(data []string) []string {
	normalized := make([]string, len(data))
	for i, item := range data {
		normalized[i] = NormalizeString(item)
	}
	return RemoveDuplicates(normalized)
}

func main() {
	rawData := []string{"  Apple", "banana", "  apple", "Banana", "Cherry  "}
	cleaned := CleanData(rawData)
	fmt.Println("Cleaned data:", cleaned)
}
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func TrimWhitespace(slice []string) []string {
	result := make([]string, len(slice))
	for i, item := range slice {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func CleanData(input []string) []string {
	trimmed := TrimWhitespace(input)
	deduped := RemoveDuplicates(trimmed)
	return deduped
}