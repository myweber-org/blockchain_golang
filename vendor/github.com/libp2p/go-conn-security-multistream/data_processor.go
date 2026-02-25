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
	Tags      []string `json:"tags"`
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func filterTags(tags []string) []string {
	var filtered []string
	seen := make(map[string]bool)
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag != "" && !seen[tag] {
			seen[tag] = true
			filtered = append(filtered, tag)
		}
	}
	return filtered
}

func processUserProfile(rawData []byte) (*UserProfile, error) {
	var profile UserProfile
	if err := json.Unmarshal(rawData, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if profile.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", profile.ID)
	}

	profile.Username = normalizeUsername(profile.Username)
	if profile.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if !validateEmail(profile.Email) {
		return nil, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", profile.Age)
	}

	profile.Tags = filterTags(profile.Tags)

	return &profile, nil
}

func main() {
	jsonData := []byte(`{
		"id": 1001,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"tags": ["golang", "backend", "GOLANG", "  ", "devops"]
	}`)

	processedProfile, err := processUserProfile(jsonData)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	fmt.Printf("Processed User Profile:\n")
	fmt.Printf("ID: %d\n", processedProfile.ID)
	fmt.Printf("Username: %s\n", processedProfile.Username)
	fmt.Printf("Email: %s\n", processedProfile.Email)
	fmt.Printf("Age: %d\n", processedProfile.Age)
	fmt.Printf("Active: %t\n", processedProfile.Active)
	fmt.Printf("Tags: %v\n", processedProfile.Tags)
}