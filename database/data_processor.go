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

func normalizeEmail(email string) (string, error) {
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

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	email, err := normalizeEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("email validation failed: %v", err)
	}
	data.Email = email

	if err := validateUsername(data.Username); err != nil {
		return nil, fmt.Errorf("username validation failed: %v", err)
	}

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email": "  TEST@Example.COM ", "username": "valid_user123", "age": 25}`
	processed, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}
package main

import "fmt"

func MovingAverage(data []float64, window int) []float64 {
    if window <= 0 || window > len(data) {
        return nil
    }

    result := make([]float64, len(data)-window+1)
    var sum float64

    for i := 0; i < window; i++ {
        sum += data[i]
    }
    result[0] = sum / float64(window)

    for i := window; i < len(data); i++ {
        sum = sum - data[i-window] + data[i]
        result[i-window+1] = sum / float64(window)
    }

    return result
}

func main() {
    values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
    averaged := MovingAverage(values, 3)
    fmt.Println("Original:", values)
    fmt.Println("Moving Average (window=3):", averaged)
}
package data_processor

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func ParseAndValidateJSON(rawData []byte, target interface{}, requiredFields []string) error {
	if err := json.Unmarshal(rawData, target); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &dataMap); err != nil {
		return fmt.Errorf("failed to parse JSON into map: %w", err)
	}

	for _, field := range requiredFields {
		if value, exists := dataMap[field]; !exists || value == nil || value == "" {
			return ValidationError{
				Field:   field,
				Message: "field is required and cannot be empty",
			}
		}
	}

	return nil
}