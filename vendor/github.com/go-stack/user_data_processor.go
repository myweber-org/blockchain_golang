
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if strings.TrimSpace(data.Username) == "" {
		errors = append(errors, "username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		errors = append(errors, "username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(data.Email) {
		errors = append(errors, "invalid email format")
	}

	if data.Age < 18 || data.Age > 120 {
		errors = append(errors, "age must be between 18 and 120")
	}

	return len(errors) == 0, errors
}

func ProcessUserJSON(jsonData string) (*UserData, error) {
	var user UserData
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	valid, errors := ValidateUserData(user)
	if !valid {
		return &user, fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return &user, nil
}

func main() {
	jsonInput := `{"username": "john_doe", "email": "john@example.com", "age": 25}`
	user, err := ProcessUserJSON(jsonInput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}