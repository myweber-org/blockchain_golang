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
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

func FormatJSON(input string) (string, error) {
    var data interface{}
    err := json.Unmarshal([]byte(input), &data)
    if err != nil {
        return "", fmt.Errorf("invalid JSON: %w", err)
    }

    formatted, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to format JSON: %w", err)
    }

    return string(formatted), nil
}

func ValidateJSON(input string) bool {
    var js json.RawMessage
    return json.Unmarshal([]byte(input), &js) == nil
}

func main() {
    sample := `{"name":"test","value":123,"active":true}`
    fmt.Println("Is valid?", ValidateJSON(sample))

    formatted, err := FormatJSON(sample)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Formatted JSON:")
    fmt.Println(formatted)
}