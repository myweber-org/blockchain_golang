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

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  JohnDoe  ","age":25}`
	processedData, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed Data: %+v\n", processedData)
}package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove any HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(trimmed, "")
	
	// Escape potentially dangerous characters
	re = regexp.MustCompile(`[<>"'&]`)
	sanitized := re.ReplaceAllStringFunc(cleaned, func(match string) string {
		switch match {
		case "<":
			return "&lt;"
		case ">":
			return "&gt;"
		case "\"":
			return "&quot;"
		case "'":
			return "&#39;"
		case "&":
			return "&amp;"
		default:
			return match
		}
	})
	
	return sanitized
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return usernameRegex.MatchString(username)
}