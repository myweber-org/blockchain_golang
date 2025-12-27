
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

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func FilterInactiveUsers(users []UserProfile) []UserProfile {
	var activeUsers []UserProfile
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func TransformUserData(users []UserProfile) ([]UserProfile, error) {
	var transformed []UserProfile
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			return nil, fmt.Errorf("invalid email for user %d: %s", user.ID, user.Email)
		}
		user.Username = NormalizeUsername(user.Username)
		if user.Age < 0 || user.Age > 150 {
			return nil, fmt.Errorf("invalid age for user %d: %d", user.ID, user.Age)
		}
		transformed = append(transformed, user)
	}
	return transformed, nil
}

func ProcessUserJSON(jsonData []byte) ([]UserProfile, error) {
	var users []UserProfile
	if err := json.Unmarshal(jsonData, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	
	transformed, err := TransformUserData(users)
	if err != nil {
		return nil, err
	}
	
	return FilterInactiveUsers(transformed), nil
}

func main() {
	jsonInput := `[
		{"id":1,"username":" JohnDoe ","email":"john@example.com","age":25,"active":true,"tags":["admin","user"]},
		{"id":2,"username":"jane_smith","email":"invalid-email","age":30,"active":false,"tags":["user"]},
		{"id":3,"username":"BOB","email":"bob@example.org","age":200,"active":true,"tags":[]}
	]`
	
	users, err := ProcessUserJSON([]byte(jsonInput))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	
	fmt.Printf("Processed %d active users\n", len(users))
	for _, user := range users {
		fmt.Printf("ID: %d, Username: %s, Email: %s\n", user.ID, user.Username, user.Email)
	}
}