package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseUser(jsonData []byte) (*User, error) {
	var user User
	err := json.Unmarshal(jsonData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	validJSON := []byte(`{"id": 123, "name": "John Doe", "email": "john@example.com"}`)
	user, err := ValidateAndParseUser(validJSON)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed user: %+v\n", user)

	invalidJSON := []byte(`{"id": -5, "name": "", "email": ""}`)
	_, err = ValidateAndParseUser(invalidJSON)
	if err != nil {
		fmt.Println("Expected error:", err)
	}
}