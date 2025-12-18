package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string, maxLength int) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)

	// Limit length
	if len(trimmed) > maxLength {
		trimmed = trimmed[:maxLength]
	}

	// Remove potentially dangerous characters
	re := regexp.MustCompile(`[<>{};]`)
	sanitized := re.ReplaceAllString(trimmed, "")

	return sanitized
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}