
package main

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func TransformPhoneNumber(phone string) (string, error) {
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	if len(cleaned) < 10 {
		return "", errors.New("phone number too short")
	}
	if len(cleaned) > 15 {
		return "", errors.New("phone number too long")
	}
	return "+" + cleaned, nil
}

func ProcessUserData(email, username, phone string) (map[string]string, error) {
	result := make(map[string]string)

	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	result["email"] = email

	normalizedUser := NormalizeUsername(username)
	result["username"] = normalizedUser

	transformedPhone, err := TransformPhoneNumber(phone)
	if err != nil {
		return nil, err
	}
	result["phone"] = transformedPhone

	return result, nil
}