
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
}package main

import (
	"fmt"
	"strings"
)

func DeduplicateStrings(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func main() {
	emails := []string{"test@example.com", "duplicate@test.org", "test@example.com", "invalid-email"}
	uniqueEmails := DeduplicateStrings(emails)
	fmt.Println("Unique emails:", uniqueEmails)

	for _, email := range uniqueEmails {
		if ValidateEmail(email) {
			fmt.Printf("%s is valid\n", email)
		} else {
			fmt.Printf("%s is invalid\n", email)
		}
	}
}