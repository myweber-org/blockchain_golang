package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("invalid email format")
	}
	return email, nil
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	email, err := normalizeEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("email validation failed: %v", err)
	}
	data.Email = email

	if err := validateUsername(data.Username); err != nil {
		return nil, fmt.Errorf("username validation failed: %v", err)
	}

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email": "  TEST@Example.COM ", "username": "valid_user123", "age": 25}`
	processed, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}