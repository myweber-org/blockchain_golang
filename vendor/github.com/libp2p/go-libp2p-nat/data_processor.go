package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(trimmed, "")
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Username = SanitizeInput(data.Username)
	data.Email = SanitizeInput(data.Email)
	data.Comments = SanitizeInput(data.Comments)

	if !ValidateEmail(data.Email) {
		return data, ErrInvalidEmail
	}

	if len(data.Username) < 3 {
		return data, ErrUsernameTooShort
	}

	return data, nil
}

var (
	ErrInvalidEmail     = NewValidationError("invalid email format")
	ErrUsernameTooShort = NewValidationError("username must be at least 3 characters")
)

type ValidationError struct {
	Message string
}

func NewValidationError(msg string) ValidationError {
	return ValidationError{Message: msg}
}

func (e ValidationError) Error() string {
	return e.Message
}