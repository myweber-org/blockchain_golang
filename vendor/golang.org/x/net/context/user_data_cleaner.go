package main

import (
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email       string
	PhoneNumber string
}

func CleanEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("invalid email format")
	}
	return email, nil
}

func FormatPhoneNumber(phone string) string {
	re := regexp.MustCompile(`\D`)
	cleaned := re.ReplaceAllString(phone, "")
	if len(cleaned) == 10 {
		return fmt.Sprintf("(%s) %s-%s", cleaned[0:3], cleaned[3:6], cleaned[6:])
	}
	return phone
}

func CleanUserData(user UserData) (UserData, error) {
	cleanedEmail, err := CleanEmail(user.Email)
	if err != nil {
		return UserData{}, err
	}
	formattedPhone := FormatPhoneNumber(user.PhoneNumber)
	return UserData{
		Email:       cleanedEmail,
		PhoneNumber: formattedPhone,
	}, nil
}

func main() {
	testUser := UserData{
		Email:       "  TEST@Example.COM  ",
		PhoneNumber: "123-456-7890",
	}
	cleaned, err := CleanUserData(testUser)
	if err != nil {
		fmt.Printf("Error cleaning data: %v\n", err)
		return
	}
	fmt.Printf("Cleaned Email: %s\n", cleaned.Email)
	fmt.Printf("Formatted Phone: %s\n", cleaned.PhoneNumber)
}