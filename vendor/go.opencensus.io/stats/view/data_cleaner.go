
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, trimmed)
	
	re := regexp.MustCompile(`\s+`)
	normalized = re.ReplaceAllString(normalized, " ")
	
	return normalized
}

func NormalizeWhitespace(input string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(input), " ")
}
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
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
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if !encountered[elements[v]] {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func main() {
	data := []string{" Apple", "banana ", "Apple", "  BANANA", "Cherry"}
	cleaned := []string{}

	for _, item := range data {
		cleaned = append(cleaned, CleanString(item))
	}

	unique := RemoveDuplicates(cleaned)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned and unique:", unique)
}