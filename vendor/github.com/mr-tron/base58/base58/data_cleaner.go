
package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes leading/trailing whitespace and collapses multiple spaces
func SanitizeInput(input string) string {
	// Trim spaces from start and end
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
}