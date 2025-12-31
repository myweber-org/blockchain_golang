
package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
}