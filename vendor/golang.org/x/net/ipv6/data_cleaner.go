package utils

import (
	"regexp"
	"strings"
)

var (
	whitespaceRegex = regexp.MustCompile(`\s+`)
	specialCharRegex = regexp.MustCompile(`[^\w\s-]`)
)

// SanitizeInput removes excessive whitespace and special characters from a string
func SanitizeInput(input string) string {
	if input == "" {
		return input
	}
	
	// Trim leading/trailing whitespace
	cleaned := strings.TrimSpace(input)
	
	// Replace multiple whitespace characters with single space
	cleaned = whitespaceRegex.ReplaceAllString(cleaned, " ")
	
	// Remove special characters except hyphens and underscores
	cleaned = specialCharRegex.ReplaceAllString(cleaned, "")
	
	return cleaned
}

// NormalizeSpacing ensures consistent spacing around punctuation
func NormalizeSpacing(text string) string {
	patterns := map[string]*regexp.Regexp{
		"spacesBeforePunctuation": regexp.MustCompile(`\s+([.,!?;:])`),
		"noSpaceAfterOpening":     regexp.MustCompile(`\(\s+`),
		"noSpaceBeforeClosing":    regexp.MustCompile(`\s+\)`),
	}
	
	result := text
	for _, re := range patterns {
		if re == patterns["spacesBeforePunctuation"] {
			result = re.ReplaceAllString(result, "$1")
		} else {
			result = re.ReplaceAllString(result, "(")
		}
	}
	
	return result
}