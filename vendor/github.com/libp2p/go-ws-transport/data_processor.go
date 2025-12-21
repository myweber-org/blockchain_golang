package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndTransform(data UserData) (UserData, error) {
	if strings.TrimSpace(data.Username) == "" {
		return UserData{}, fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return UserData{}, fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := UserData{
		Username: strings.ToLower(strings.TrimSpace(data.Username)),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}
	return transformed, nil
}

func main() {
	user := UserData{
		Username: "  TestUser  ",
		Email:    "EXAMPLE@DOMAIN.COM",
		Age:      25,
	}

	result, err := ValidateAndTransform(user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed: %+v\n", result)
}