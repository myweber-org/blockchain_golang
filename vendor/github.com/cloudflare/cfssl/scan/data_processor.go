
package main

import (
    "errors"
    "regexp"
    "strings"
)

type UserData struct {
    Email    string
    Username string
    Age      int
}

func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("email cannot be empty")
    }
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
    }
    return nil
}

func SanitizeUsername(username string) string {
    username = strings.TrimSpace(username)
    username = strings.ToLower(username)
    return username
}

func ProcessUserData(data UserData) (UserData, error) {
    if err := ValidateEmail(data.Email); err != nil {
        return UserData{}, err
    }

    sanitizedUsername := SanitizeUsername(data.Username)
    if sanitizedUsername == "" {
        return UserData{}, errors.New("username cannot be empty after sanitization")
    }

    if data.Age < 0 || data.Age > 150 {
        return UserData{}, errors.New("age must be between 0 and 150")
    }

    return UserData{
        Email:    data.Email,
        Username: sanitizedUsername,
        Age:      data.Age,
    }, nil
}package main

import (
	"errors"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}
	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func TransformUserName(data *UserData) {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
}

func ProcessUserInput(data UserData) (UserData, error) {
	if err := ValidateUserData(data); err != nil {
		return UserData{}, err
	}
	TransformUserName(&data)
	return data, nil
}