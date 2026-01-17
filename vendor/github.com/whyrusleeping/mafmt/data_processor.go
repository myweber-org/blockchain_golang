package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func transformData(input []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(input, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	data.Username = sanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 120 {
		return nil, fmt.Errorf("age must be between 0 and 120")
	}

	return &data, nil
}

func main() {
	jsonInput := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processedData, err := transformData([]byte(jsonInput))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed Data: %+v\n", processedData)
}