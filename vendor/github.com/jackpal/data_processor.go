
package main

import (
	"regexp"
	"strings"
	"unicode"
)

func CleanInput(input string) string {
	processed := strings.TrimSpace(input)
	processed = collapseWhitespace(processed)
	processed = removeSpecialChars(processed)
	return strings.ToLower(processed)
}

func collapseWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}

func removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}