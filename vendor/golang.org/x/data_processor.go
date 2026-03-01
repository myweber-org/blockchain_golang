package main

import (
	"errors"
	"fmt"
)

type Data struct {
	ID    int
	Value string
}

func ProcessData(data []Data) (map[int]string, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data slice provided")
	}

	processed := make(map[int]string)
	for _, item := range data {
		if item.ID <= 0 {
			return nil, fmt.Errorf("invalid ID %d found in data", item.ID)
		}
		if item.Value == "" {
			return nil, fmt.Errorf("empty value for ID %d", item.ID)
		}
		processed[item.ID] = item.Value + "_processed"
	}
	return processed, nil
}

func ValidateData(data Data) error {
	if data.ID <= 0 {
		return errors.New("ID must be positive integer")
	}
	if data.Value == "" {
		return errors.New("value cannot be empty")
	}
	if len(data.Value) > 100 {
		return errors.New("value exceeds maximum length of 100 characters")
	}
	return nil
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
	
	// Trim spaces from beginning and end
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase for normalization
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
}