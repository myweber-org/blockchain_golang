package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	normalizeCase bool
	trimSpaces    bool
}

func NewDataCleaner(options ...func(*DataCleaner)) *DataCleaner {
	dc := &DataCleaner{
		normalizeCase: true,
		trimSpaces:    true,
	}
	for _, option := range options {
		option(dc)
	}
	return dc
}

func WithCaseNormalization(enabled bool) func(*DataCleaner) {
	return func(dc *DataCleaner) {
		dc.normalizeCase = enabled
	}
}

func WithSpaceTrimming(enabled bool) func(*DataCleaner) {
	return func(dc *DataCleaner) {
		dc.trimSpaces = enabled
	}
}

func (dc *DataCleaner) CleanString(input string) string {
	result := input

	if dc.trimSpaces {
		result = strings.TrimSpace(result)
	}

	if dc.normalizeCase {
		result = strings.ToLower(result)
	}

	return result
}

func (dc *DataCleaner) DeduplicateStrings(items []string) []string {
	seen := make(map[string]struct{})
	var unique []string

	for _, item := range items {
		cleaned := dc.CleanString(item)
		if _, exists := seen[cleaned]; !exists {
			seen[cleaned] = struct{}{}
			unique = append(unique, cleaned)
		}
	}
	return unique
}

func NormalizeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	localPart := strings.Split(parts[0], "+")[0]
	localPart = strings.ReplaceAll(localPart, ".", "")
	return strings.ToLower(localPart + "@" + parts[1])
}

func main() {
	cleaner := NewDataCleaner()

	data := []string{
		"  Apple  ",
		"apple",
		"APPLE",
		"Banana",
		"  banana ",
	}

	cleaned := cleaner.DeduplicateStrings(data)
	fmt.Println("Deduplicated data:", cleaned)

	emails := []string{
		"John.Doe+test@example.com",
		"johndoe@EXAMPLE.COM",
		"john.doe@example.com",
	}

	fmt.Println("Normalized emails:")
	for _, email := range emails {
		fmt.Println(NormalizeEmail(email))
	}
}