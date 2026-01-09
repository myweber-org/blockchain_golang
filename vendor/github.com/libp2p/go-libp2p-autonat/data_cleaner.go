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