
package main

import (
    "regexp"
    "strings"
)

func SanitizeUsername(input string) string {
    re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
    sanitized := re.ReplaceAllString(input, "")
    return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(pattern)
    return re.MatchString(email)
}

func SanitizeText(input string, maxLength int) string {
    trimmed := strings.TrimSpace(input)
    if len(trimmed) > maxLength {
        trimmed = trimmed[:maxLength]
    }
    return trimmed
}
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ContainsOnlyAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}