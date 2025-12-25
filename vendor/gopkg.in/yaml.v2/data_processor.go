
package data_processor

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	pattern := `^[a-zA-Z0-9\s\.\-_@]+$`
	matched, err := regexp.MatchString(pattern, trimmed)
	if err != nil || !matched {
		return "", false
	}

	return trimmed, true
}