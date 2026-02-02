
package utils

import (
	"regexp"
	"strings"
)

var (
	whitespaceRegex = regexp.MustCompile(`\s+`)
	htmlTagRegex    = regexp.MustCompile(`<[^>]*>`)
)

func SanitizeInput(input string) string {
	if input == "" {
		return input
	}

	cleaned := strings.TrimSpace(input)
	cleaned = htmlTagRegex.ReplaceAllString(cleaned, "")
	cleaned = whitespaceRegex.ReplaceAllString(cleaned, " ")
	
	return cleaned
}

func NormalizeSpaces(input string) string {
	return whitespaceRegex.ReplaceAllString(strings.TrimSpace(input), " ")
}