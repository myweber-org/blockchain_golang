
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, trimmed)
	
	re := regexp.MustCompile(`\s+`)
	normalized = re.ReplaceAllString(normalized, " ")
	
	return normalized
}

func NormalizeWhitespace(input string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(input), " ")
}