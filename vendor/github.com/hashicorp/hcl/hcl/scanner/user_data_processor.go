
package main

import (
    "regexp"
    "strings"
)

type UserData struct {
    Username string
    Email    string
    Age      int
}

func ValidateUsername(username string) bool {
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    return validUsername.MatchString(username)
}

func SanitizeEmail(email string) string {
    trimmedEmail := strings.TrimSpace(email)
    return strings.ToLower(trimmedEmail)
}

func ValidateUserAge(age int) bool {
    return age >= 18 && age <= 120
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
    if !ValidateUsername(username) {
        return UserData{}, ErrInvalidUsername
    }

    sanitizedEmail := SanitizeEmail(email)
    if !ValidateUserAge(age) {
        return UserData{}, ErrInvalidAge
    }

    return UserData{
        Username: username,
        Email:    sanitizedEmail,
        Age:      age,
    }, nil
}

var (
    ErrInvalidUsername = errors.New("invalid username format")
    ErrInvalidAge      = errors.New("age must be between 18 and 120")
)