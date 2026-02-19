package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Remove any non-printable characters
	clean := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, input)

	// Replace multiple whitespaces with single space
	re := regexp.MustCompile(`\s+`)
	clean = re.ReplaceAllString(clean, " ")

	// Trim leading and trailing whitespace
	clean = strings.TrimSpace(clean)

	return clean
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func RemoveExtraSpaces(input string) string {
	re := regexp.MustCompile(`\s{2,}`)
	return re.ReplaceAllString(input, " ")
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	normalizeCase bool
	trimSpaces    bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		normalizeCase: true,
		trimSpaces:    true,
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	result := input

	if dc.trimSpaces {
		result = strings.TrimSpace(result)
	}

	if dc.normalizeCase {
		result = strings.ToLower(result)
	}

	return result
}

func (dc *DataCleaner) DeduplicateStrings(strings []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, str := range strings {
		normalized := dc.NormalizeString(str)
		if _, exists := seen[normalized]; !exists {
			seen[normalized] = struct{}{}
			result = append(result, normalized)
		}
	}

	return result
}

func main() {
	cleaner := NewDataCleaner()

	sampleData := []string{
		"  Apple  ",
		"apple",
		" BANANA",
		"banana ",
		"Cherry",
		"CHERRY",
	}

	fmt.Println("Original data:", sampleData)

	cleaned := cleaner.DeduplicateStrings(sampleData)
	fmt.Println("Cleaned data:", cleaned)
}