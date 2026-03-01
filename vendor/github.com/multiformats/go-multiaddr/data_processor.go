
package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndNormalize(data *UserData) error {
	if data == nil {
		return errors.New("user data cannot be nil")
	}

	data.Username = strings.TrimSpace(data.Username)
	if data.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	data.Email = strings.TrimSpace(data.Email)
	if data.Email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 {
		return errors.New("age cannot be negative")
	}
	if data.Age > 150 {
		return errors.New("age is unrealistically high")
	}

	return nil
}