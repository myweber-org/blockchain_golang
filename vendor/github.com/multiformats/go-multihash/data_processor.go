package main

import (
	"encoding/json"
	"errors"
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
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 120 {
		return errors.New("age must be between 0 and 120")
	}

	return nil
}

func TransformUsername(profile *UserProfile) {
	profile.Username = strings.ToLower(strings.TrimSpace(profile.Username))
}

func ProcessUserData(jsonData []byte) (*UserProfile, error) {
	var profile UserProfile
	err := json.Unmarshal(jsonData, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	TransformUsername(&profile)

	err = ValidateUserProfile(profile)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return &profile, nil
}

func main() {
	sampleData := `{
		"id": 123,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true
	}`

	profile, err := ProcessUserData([]byte(sampleData))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed user profile: %+v\n", profile)
}