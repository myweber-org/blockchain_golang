package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", profile.ID)
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 120 {
		return fmt.Errorf("age must be between 0 and 120")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(transformed.Username)
	transformed.Email = strings.ToLower(transformed.Email)
	transformed.Active = true

	if transformed.CreatedAt.IsZero() {
		transformed.CreatedAt = time.Now().UTC()
	}

	return transformed
}

func ProcessUserProfile(data []byte) (UserProfile, error) {
	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return UserProfile{}, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, fmt.Errorf("validation failed: %w", err)
	}

	transformedProfile := TransformProfile(profile)
	return transformedProfile, nil
}

func main() {
	jsonData := []byte(`{
		"id": 123,
		"username": "JohnDoe",
		"email": "JOHN@EXAMPLE.COM",
		"age": 30,
		"created_at": "2023-01-15T10:30:00Z",
		"active": false
	}`)

	processedProfile, err := ProcessUserProfile(jsonData)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	output, _ := json.MarshalIndent(processedProfile, "", "  ")
	fmt.Println("Processed user profile:")
	fmt.Println(string(output))
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	jsonData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25,"active":true}`)
	processedData, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}