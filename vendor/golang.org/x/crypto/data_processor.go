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
	if data.Username == "" {
		return UserData{}, fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return UserData{}, fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := UserData{
		Username: strings.TrimSpace(data.Username),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}
	return transformed, nil
}

func main() {
	sampleData := UserData{
		Username: "  JohnDoe  ",
		Email:    "John@Example.COM",
		Age:      30,
	}

	result, err := ValidateAndTransform(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed Data: %+v\n", result)
}