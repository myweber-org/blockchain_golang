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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Username = SanitizeInput(data.Username)
	data.Email = SanitizeInput(data.Email)
	data.Comments = SanitizeInput(data.Comments)

	if !ValidateEmail(data.Email) {
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