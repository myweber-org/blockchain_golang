
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
package main

import (
    "errors"
    "strings"
)

type UserData struct {
    Name  string
    Email string
    Age   int
}

func ValidateUserData(data UserData) error {
    if strings.TrimSpace(data.Name) == "" {
        return errors.New("name cannot be empty")
    }
    if !strings.Contains(data.Email, "@") {
        return errors.New("invalid email format")
    }
    if data.Age < 0 || data.Age > 150 {
        return errors.New("age must be between 0 and 150")
    }
    return nil
}

func TransformUserData(data UserData) UserData {
    return UserData{
        Name:  strings.ToUpper(strings.TrimSpace(data.Name)),
        Email: strings.ToLower(strings.TrimSpace(data.Email)),
        Age:   data.Age,
    }
}