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
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Deduplicate(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, normalized)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && len(item) < 100
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"apple",
		"  APPLE  ",
		"banana",
		"",
		"banana",
		"orange",
		"Orange",
		strings.Repeat("x", 150),
	}
	
	cleaned := cleaner.Deduplicate(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	fmt.Printf("Unique count: %d\n", len(cleaned))
	
	cleaner.Reset()
	secondBatch := []string{"apple", "grape"}
	secondCleaned := cleaner.Deduplicate(secondBatch)
	fmt.Printf("Second batch: %v\n", secondCleaned)
}