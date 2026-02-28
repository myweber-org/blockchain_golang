package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data *UserData) {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	user := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	TransformUsername(&user)

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	return user, nil
}

func main() {
	user, err := ProcessUserInput("  JohnDoe  ", "john@example.com", 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}
package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !validPattern.MatchString(trimmed) {
		return "", false
	}
	return trimmed, true
}

func ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

func ProcessUserData(username, email string) (map[string]string, error) {
	sanitizedUsername, valid := SanitizeUsername(username)
	if !valid {
		return nil, &InvalidDataError{Field: "username", Value: username}
	}

	if !ValidateEmail(email) {
		return nil, &InvalidDataError{Field: "email", Value: email}
	}

	return map[string]string{
		"username": sanitizedUsername,
		"email":    strings.ToLower(strings.TrimSpace(email)),
	}, nil
}

type InvalidDataError struct {
	Field string
	Value string
}

func (e *InvalidDataError) Error() string {
	return "invalid data for field: " + e.Field
}