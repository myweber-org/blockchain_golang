package data_processor

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeString(input string) string {
	cleaned := CleanInput(input)
	
	// Remove special characters except alphanumeric and spaces
	re := regexp.MustCompile(`[^a-z0-9\s]`)
	normalized := re.ReplaceAllString(cleaned, "")
	
	return normalized
}

func ValidateInput(input string, minLength int) bool {
	if len(input) < minLength {
		return false
	}
	
	// Check if input contains at least one letter
	re := regexp.MustCompile(`[a-z]`)
	return re.MatchString(input)
}