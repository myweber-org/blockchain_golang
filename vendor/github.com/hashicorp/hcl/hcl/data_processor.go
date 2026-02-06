
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func NormalizeName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

func TransformUserData(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	if !ValidateEmail(user.Email) {
		return nil, fmt.Errorf("invalid email format: %s", user.Email)
	}

	user.Name = NormalizeName(user.Name)

	if user.Age < 0 || user.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", user.Age)
	}

	return &user, nil
}

func ProcessDataBatch(rawDataList [][]byte) ([]UserData, []error) {
	var validUsers []UserData
	var errors []error

	for i, rawData := range rawDataList {
		user, err := TransformUserData(rawData)
		if err != nil {
			errors = append(errors, fmt.Errorf("item %d: %w", i, err))
			continue
		}
		validUsers = append(validUsers, *user)
	}

	return validUsers, errors
}