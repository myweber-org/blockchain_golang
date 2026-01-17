
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