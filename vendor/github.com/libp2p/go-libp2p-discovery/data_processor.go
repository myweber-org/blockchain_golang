
package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	for _, ch := range username {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Email = NormalizeEmail(data.Email)
	if !ValidateUsername(data.Username) {
		return data, fmt.Errorf("invalid username format")
	}
	if data.Age < 0 || data.Age > 150 {
		return data, fmt.Errorf("age out of valid range")
	}
	return data, nil
}

func main() {
	user := UserData{
		Email:    "  TEST@Example.COM ",
		Username: "valid_user_123",
		Age:      25,
	}
	processed, err := ProcessUserData(user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}