
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Remove any non-printable characters
	clean := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, input)

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	clean = spaceRegex.ReplaceAllString(clean, " ")

	// Trim leading and trailing whitespace
	clean = strings.TrimSpace(clean)

	return clean
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func RemoveSpecialChars(input string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return reg.ReplaceAllString(input, "")
}