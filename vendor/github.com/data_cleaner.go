package main

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
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)

	// Remove extra internal whitespace
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(trimmed, " ")

	// Remove non-printable characters
	var result strings.Builder
	for _, r := range normalized {
		if unicode.IsPrint(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; !exists {
			seen[trimmed] = struct{}{}
			result = append(result, trimmed)
		}
	}
	return result
}

func (dc *DataCleaner) Clean() []string {
	return dc.RemoveDuplicates()
}

func main() {
	rawData := []string{"  apple ", "banana", "  apple", "cherry  ", "", "banana", "date"}
	cleaner := NewDataCleaner(rawData)
	cleaned := cleaner.Clean()
	fmt.Println("Cleaned data:", cleaned)
}