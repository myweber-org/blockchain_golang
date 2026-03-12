
package data_processor

import (
	"regexp"
	"strings"
)

func CleanInput(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}

	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?-]+$`)
	if !validPattern.MatchString(trimmed) {
		return "", ErrInvalidCharacters
	}

	cleaned := regexp.MustCompile(`\s+`).ReplaceAllString(trimmed, " ")
	return cleaned, nil
}

var (
	ErrEmptyInput        = errors.New("input string cannot be empty")
	ErrInvalidCharacters = errors.New("input contains invalid characters")
)