package datautils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&;]`)
	sanitized := re.ReplaceAllString(trimmed, "")
	return sanitized
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}