package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes leading/trailing whitespace, reduces multiple spaces to single,
// and removes any non-printable characters from the input string.
func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned := spaceRegex.ReplaceAllString(trimmed, " ")
	
	// Remove non-printable characters
	printableRegex := regexp.MustCompile(`[^[:print:]]`)
	final := printableRegex.ReplaceAllString(cleaned, "")
	
	return final
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
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

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}