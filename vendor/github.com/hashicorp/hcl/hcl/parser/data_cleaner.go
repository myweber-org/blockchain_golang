package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces from beginning and end
	cleaned = strings.TrimSpace(cleaned)
	
	// Remove non-printable characters
	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, cleaned)
	
	// Convert to lowercase for consistency
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeWhitespace(input string) string {
	// Replace all whitespace variations with single space
	re := regexp.MustCompile(`[\s\p{Zs}]+`)
	return re.ReplaceAllString(input, " ")
}