
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

func sanitizeInput(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func processUserData(data UserData) (UserData, error) {
	data.Username = sanitizeInput(data.Username)
	data.Comments = sanitizeInput(data.Comments)

	if !validateEmail(data.Email) {
		return data, &ValidationError{Field: "email", Message: "invalid email format"}
	}

	if len(data.Username) < 3 || len(data.Username) > 50 {
		return data, &ValidationError{Field: "username", Message: "username must be between 3 and 50 characters"}
	}

	return data, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}