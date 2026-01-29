
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	Score     float64 `json:"score"`
}

func ValidateJSON(raw []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if data.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if data.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if data.Score < 0 || data.Score > 100 {
		return nil, fmt.Errorf("score must be between 0 and 100")
	}

	return &data, nil
}

func ProcessUserData(jsonStr string) {
	userData, err := ValidateJSON([]byte(jsonStr))
	if err != nil {
		log.Printf("Validation error: %v", err)
		return
	}

	fmt.Printf("User ID: %d\n", userData.ID)
	fmt.Printf("Username: %s\n", userData.Username)
	fmt.Printf("Email: %s\n", userData.Email)
	fmt.Printf("Active: %t\n", userData.Active)
	fmt.Printf("Score: %.2f\n", userData.Score)
}

func main() {
	sampleJSON := `{
		"id": 42,
		"username": "john_doe",
		"email": "john@example.com",
		"active": true,
		"score": 85.5
	}`

	ProcessUserData(sampleJSON)
}