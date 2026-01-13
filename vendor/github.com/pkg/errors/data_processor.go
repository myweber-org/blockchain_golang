package main

import (
	"errors"
	"strings"
)

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username must not exceed 50 characters")
	}
	if strings.ContainsAny(username, "!@#$%^&*()") {
		return errors.New("username contains invalid special characters")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformUserData(username, email string) (string, string, error) {
	if err := ValidateUsername(username); err != nil {
		return "", "", err
	}
	normalizedEmail := NormalizeEmail(email)
	return username, normalizedEmail, nil
}