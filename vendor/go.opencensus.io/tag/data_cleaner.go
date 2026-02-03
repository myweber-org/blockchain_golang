
package main

import (
	"strings"
	"unicode"
)

func SanitizeCSVField(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if r == '"' {
			builder.WriteRune('"')
			builder.WriteRune('"')
		} else if r == ',' || r == '\n' || r == '\r' {
			builder.WriteRune(' ')
		} else if unicode.IsGraphic(r) && !unicode.IsControl(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}
	result := builder.String()
	if strings.ContainsAny(result, ",\"\n\r") {
		return "\"" + result + "\""
	}
	return result
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
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}