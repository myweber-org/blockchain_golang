package data

import (
	"strings"
	"unicode"
)

// CleanString removes extra whitespace and normalizes line endings
func CleanString(input string) string {
	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces/newlines with single space
	var result strings.Builder
	result.Grow(len(trimmed))
	
	prevSpace := false
	for _, r := range trimmed {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}
	
	return result.String()
}

// NormalizeWhitespace converts all whitespace characters to single spaces
func NormalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// IsASCII checks if string contains only ASCII characters
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}