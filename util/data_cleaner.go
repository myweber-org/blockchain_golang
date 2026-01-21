package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput cleans user-provided strings by trimming whitespace
// and removing potentially dangerous special characters
func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove any non-alphanumeric characters except spaces and basic punctuation
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s.,!?-]`)
	cleaned := reg.ReplaceAllString(trimmed, "")
	
	// Replace multiple spaces with single space
	spaceReg := regexp.MustCompile(`\s+`)
	final := spaceReg.ReplaceAllString(cleaned, " ")
	
	return final
}