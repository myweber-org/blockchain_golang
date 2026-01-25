package main

import (
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}

func SanitizeEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	return strings.ToLower(trimmed)
}

func ProcessUserInput(username, email string) (*User, error) {
	if !ValidateUsername(username) {
		return nil, &ValidationError{"Invalid username format"}
	}

	cleanEmail := SanitizeEmail(email)
	user := &User{
		Username: username,
		Email:    cleanEmail,
	}
	return user, nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}