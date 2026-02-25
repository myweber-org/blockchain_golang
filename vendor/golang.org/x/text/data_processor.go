
package data_processor

import (
	"regexp"
	"strings"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dc.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return strings.ToLower(normalized)
}

func (dc *DataCleaner) RemoveSpecialChars(input string, keepPattern string) string {
	if keepPattern == "" {
		keepPattern = `[^a-zA-Z0-9\s]`
	}
	regex := regexp.MustCompile(keepPattern)
	return regex.ReplaceAllString(input, "")
}

func (dc *DataCleaner) Tokenize(input string, delimiter string) []string {
	if delimiter == "" {
		delimiter = " "
	}
	normalized := dc.NormalizeString(input)
	return strings.Split(normalized, delimiter)
}
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces from start and end
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase for consistency
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeString(input string) string {
	cleaned := CleanInput(input)
	
	// Remove special characters except alphanumeric and spaces
	re := regexp.MustCompile(`[^a-z0-9\s]`)
	normalized := re.ReplaceAllString(cleaned, "")
	
	return normalized
}

func ProcessData(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		processed := NormalizeString(input)
		if processed != "" {
			results = append(results, processed)
		}
	}
	return results
}package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	cleaned := strings.TrimSpace(input)
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	return cleaned
}

func NormalizeCase(input string) string {
	return strings.ToLower(input)
}

func RemoveSpecialChars(input string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return reg.ReplaceAllString(input, "")
}
package data_processor

import (
	"regexp"
	"strings"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dc *DataCleaner) NormalizeWhitespace(input string) string {
	trimmed := strings.TrimSpace(input)
	return dc.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	return dc.emailRegex.MatchString(email)
}

func (dc *DataCleaner) ExtractDomain(email string) (string, bool) {
	if !dc.ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}

func (dc *DataCleaner) SanitizeInput(input string, maxLength int) string {
	normalized := dc.NormalizeWhitespace(input)
	if len(normalized) > maxLength {
		return normalized[:maxLength]
	}
	return normalized
}