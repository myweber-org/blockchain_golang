package main

import (
	"errors"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateEmail(email string) error {
	email = NormalizeEmail(email)
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	if strings.Count(email, "@") != 1 {
		return errors.New("invalid email format")
	}
	parts := strings.Split(email, "@")
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return errors.New("invalid email format")
	}
	if !strings.Contains(parts[1], ".") {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserData(username, email string, age int) (*UserData, error) {
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := ValidateAge(age); err != nil {
		return nil, err
	}
	return &UserData{
		Username: username,
		Email:    NormalizeEmail(email),
		Age:      age,
	}, nil
}