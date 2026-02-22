package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
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
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range data.Username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}

	data.Email = strings.ToLower(strings.TrimSpace(data.Email))
	if !strings.Contains(data.Email, "@") || !strings.Contains(data.Email, ".") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func ProcessUserInput(username, email string, age int) (*UserData, error) {
	user := &UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	if err := ValidateAndNormalize(user); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return user, nil
}

func main() {
	user, err := ProcessUserInput("  John_Doe-123  ", "JOHN@EXAMPLE.COM", 25)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}