
package main

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserInput(email, username string, age int) (*UserData, error) {
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	sanitizedUsername := SanitizeUsername(username)
	if sanitizedUsername == "" {
		return nil, errors.New("username cannot be empty")
	}

	if err := ValidateAge(age); err != nil {
		return nil, err
	}

	return &UserData{
		Email:    email,
		Username: sanitizedUsername,
		Age:      age,
	}, nil
}