package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes excessive whitespace and trims leading/trailing spaces from a string.
func SanitizeInput(input string) string {
	// Replace multiple spaces, tabs, and newlines with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	// Trim spaces from both ends
	return strings.TrimSpace(cleaned)
}