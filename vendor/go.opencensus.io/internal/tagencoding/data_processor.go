
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
	re := regexp.MustCompile(`<[^>]*>`)
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
		return data, &InvalidEmailError{Email: data.Email}
	}

	if len(data.Username) < 3 {
		return data, &InvalidUsernameError{Username: data.Username}
	}

	return data, nil
}

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return "Invalid email format: " + e.Email
}

type InvalidUsernameError struct {
	Username string
}

func (e *InvalidUsernameError) Error() string {
	return "Username must be at least 3 characters: " + e.Username
}package main

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

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if len(data.Username) < 3 || len(data.Username) > 20 {
		errors = append(errors, "Username must be between 3 and 20 characters")
	}

	if !emailRegex.MatchString(data.Email) {
		errors = append(errors, "Invalid email format")
	}

	if len(data.Comments) > 500 {
		errors = append(errors, "Comments must not exceed 500 characters")
	}

	data.Username = SanitizeInput(data.Username)
	data.Email = SanitizeInput(data.Email)
	data.Comments = SanitizeInput(data.Comments)

	return len(errors) == 0, errors
}

func ProcessUserData(data UserData) (UserData, []string) {
	isValid, validationErrors := ValidateUserData(data)
	if !isValid {
		return data, validationErrors
	}

	return data, nil
}