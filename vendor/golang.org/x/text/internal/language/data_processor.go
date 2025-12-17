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

func ValidateUser(data UserData) (bool, []string) {
	var errors []string

	if strings.TrimSpace(data.Username) == "" {
		errors = append(errors, "username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		errors = append(errors, "invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		errors = append(errors, "age must be between 0 and 150")
	}

	return len(errors) == 0, errors
}

func TransformUsername(data *UserData) {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(username, email string, age int) (UserData, []string) {
	user := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	TransformUsername(&user)
	valid, errors := ValidateUser(user)

	if !valid {
		return UserData{}, errors
	}

	return user, nil
}

func main() {
	user, errs := ProcessUserInput("  JohnDoe  ", "john@example.com", 25)
	if errs != nil {
		fmt.Println("Validation errors:", errs)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}