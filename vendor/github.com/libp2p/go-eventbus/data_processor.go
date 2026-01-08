package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", profile.ID)
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 120 {
		return fmt.Errorf("age must be between 0 and 120")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(transformed.Username)
	transformed.Email = strings.ToLower(transformed.Email)
	transformed.Active = true

	if transformed.CreatedAt.IsZero() {
		transformed.CreatedAt = time.Now().UTC()
	}

	return transformed
}

func ProcessUserProfile(data []byte) (UserProfile, error) {
	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return UserProfile{}, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, fmt.Errorf("validation failed: %w", err)
	}

	transformedProfile := TransformProfile(profile)
	return transformedProfile, nil
}

func main() {
	jsonData := []byte(`{
		"id": 123,
		"username": "JohnDoe",
		"email": "JOHN@EXAMPLE.COM",
		"age": 30,
		"created_at": "2023-01-15T10:30:00Z",
		"active": false
	}`)

	processedProfile, err := ProcessUserProfile(jsonData)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	output, _ := json.MarshalIndent(processedProfile, "", "  ")
	fmt.Println("Processed user profile:")
	fmt.Println(string(output))
}