
package sanitizer

import (
	"regexp"
	"strings"
)

var (
	whitespaceRegex = regexp.MustCompile(`\s+`)
	htmlTagRegex    = regexp.MustCompile(`<[^>]*>`)
)

func CleanInput(input string) string {
	cleaned := strings.TrimSpace(input)
	cleaned = htmlTagRegex.ReplaceAllString(cleaned, "")
	cleaned = whitespaceRegex.ReplaceAllString(cleaned, " ")
	return cleaned
}

func NormalizeWhitespace(input string) string {
	return whitespaceRegex.ReplaceAllString(strings.TrimSpace(input), " ")
}

func RemoveHTMLTags(input string) string {
	return htmlTagRegex.ReplaceAllString(input, "")
}