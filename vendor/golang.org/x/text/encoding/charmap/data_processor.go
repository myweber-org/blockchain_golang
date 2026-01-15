package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func SanitizeInput(input string) string {
	var result strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Username = NormalizeUsername(data.Username)
	data.Username = SanitizeInput(data.Username)

	if !ValidateEmail(data.Email) {
		return data, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return data, fmt.Errorf("age out of valid range")
	}

	return data, nil
}

func main() {
	sampleData := UserData{
		Username: "  John_Doe123  ",
		Email:    "john@example.com",
		Age:      30,
	}

	processed, err := ProcessUserData(sampleData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed user: %+v\n", processed)
}