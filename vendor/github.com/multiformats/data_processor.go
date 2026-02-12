
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?-]+$`)
	return &DataProcessor{allowedPattern: pattern}
}

func (dp *DataProcessor) SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	if !dp.allowedPattern.MatchString(trimmed) {
		return "", false
	}

	return trimmed, true
}

func (dp *DataProcessor) ProcessUserData(rawData string) (string, error) {
	sanitized, valid := dp.SanitizeInput(rawData)
	if !valid {
		return "", &InvalidInputError{Input: rawData}
	}

	processed := strings.ToUpper(sanitized)
	return processed, nil
}

type InvalidInputError struct {
	Input string
}

func (e *InvalidInputError) Error() string {
	return "input contains invalid characters or is empty"
}