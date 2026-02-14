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
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func ContainsSQLInjection(input string) bool {
	keywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "UNION", "--"}
	upperInput := strings.ToUpper(input)
	for _, keyword := range keywords {
		if strings.Contains(upperInput, keyword) {
			return true
		}
	}
	return false
}