
package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return re.ReplaceAllString(input, "")
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func ContainsProfanity(text string, blacklist []string) bool {
	lowerText := strings.ToLower(text)
	for _, word := range blacklist {
		if strings.Contains(lowerText, strings.ToLower(word)) {
			return true
		}
	}
	return false
}