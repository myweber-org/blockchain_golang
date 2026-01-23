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

	// Replace multiple whitespaces with single space
	re := regexp.MustCompile(`\s+`)
	clean = re.ReplaceAllString(clean, " ")

	// Trim leading and trailing whitespace
	clean = strings.TrimSpace(clean)

	return clean
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func RemoveExtraSpaces(input string) string {
	re := regexp.MustCompile(`\s{2,}`)
	return re.ReplaceAllString(input, " ")
}