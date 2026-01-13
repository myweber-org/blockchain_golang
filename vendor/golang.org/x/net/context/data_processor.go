package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func ParseAndValidateJSON(rawData []byte, target interface{}) error {
	if len(rawData) == 0 {
		return ValidationError{Field: "data", Message: "empty input"}
	}

	if !json.Valid(rawData) {
		return ValidationError{Field: "data", Message: "invalid JSON format"}
	}

	if err := json.Unmarshal(rawData, target); err != nil {
		if strings.Contains(err.Error(), "json: cannot unmarshal") {
			parts := strings.Split(err.Error(), " into ")
			if len(parts) > 0 {
				fieldPart := strings.TrimPrefix(parts[0], "json: cannot unmarshal ")
				return ValidationError{Field: fieldPart, Message: "type mismatch"}
			}
		}
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	validJSON := []byte(`{"id": 1, "name": "Alice", "email": "alice@example.com"}`)
	var user User

	if err := ParseAndValidateJSON(validJSON, &user); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Parsed user: %+v\n", user)
	}

	invalidJSON := []byte(`{"id": "not_a_number", "name": "Bob"}`)
	var anotherUser User
	if err := ParseAndValidateJSON(invalidJSON, &anotherUser); err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}