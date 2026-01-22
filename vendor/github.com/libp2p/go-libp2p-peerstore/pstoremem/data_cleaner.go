
package main

import (
	"fmt"
	"strings"
)

func DeduplicateStrings(slice []string) []string {
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

func NormalizeWhitespace(input string) string {
	words := strings.Fields(input)
	return strings.Join(words, " ")
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana"}
	unique := DeduplicateStrings(data)
	fmt.Println("Deduplicated:", unique)

	text := "  This   is   a  test   string.  "
	normalized := NormalizeWhitespace(text)
	fmt.Println("Normalized:", normalized)
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
	data := []int{7, 2, 5, 2, 8, 7, 1, 9, 5}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}