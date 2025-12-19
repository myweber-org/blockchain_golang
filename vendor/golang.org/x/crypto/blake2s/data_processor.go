
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Biography string `json:"biography,omitempty"`
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

func TransformProfile(rawData []byte) (*UserProfile, error) {
	var profile UserProfile
	err := json.Unmarshal(rawData, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	profile.Username = SanitizeUsername(profile.Username)

	if !ValidateEmail(profile.Email) {
		return nil, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", profile.Age)
	}

	return &profile, nil
}

func ProcessUserData(data []byte) (string, error) {
	profile, err := TransformProfile(data)
	if err != nil {
		return "", err
	}

	output, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile: %w", err)
	}

	return string(output), nil
}

func main() {
	rawJSON := `{
		"id": 42,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"biography": "Software developer"
	}`

	result, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Processed profile:")
	fmt.Println(result)
}