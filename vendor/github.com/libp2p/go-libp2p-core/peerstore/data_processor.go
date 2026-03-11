
package main

import (
    "fmt"
    "strings"
    "unicode"
)

type UserData struct {
    Username string
    Email    string
}

func NormalizeUsername(username string) string {
    trimmed := strings.TrimSpace(username)
    var result strings.Builder
    for _, r := range trimmed {
        if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
            result.WriteRune(unicode.ToLower(r))
        }
    }
    return result.String()
}

func ValidateEmail(email string) bool {
    trimmed := strings.TrimSpace(email)
    if !strings.Contains(trimmed, "@") || !strings.Contains(trimmed, ".") {
        return false
    }
    return len(trimmed) > 5 && len(trimmed) < 255
}

func ProcessUserInput(username, email string) (*UserData, error) {
    normalizedUsername := NormalizeUsername(username)
    if normalizedUsername == "" {
        return nil, fmt.Errorf("username cannot be empty after normalization")
    }

    if !ValidateEmail(email) {
        return nil, fmt.Errorf("invalid email format")
    }

    return &UserData{
        Username: normalizedUsername,
        Email:    strings.ToLower(strings.TrimSpace(email)),
    }, nil
}

func main() {
    user, err := ProcessUserInput("  John_Doe-123  ", "john@example.com")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Processed user: %+v\n", user)
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// DataPayload represents a generic structure for incoming JSON data
type DataPayload struct {
    ID    string `json:"id"`
    Value int    `json:"value"`
    Tags  []string `json:"tags"`
}

// ParseAndValidateJSON parses a JSON string into a DataPayload and performs basic validation
func ParseAndValidateJSON(input string) (*DataPayload, error) {
    if strings.TrimSpace(input) == "" {
        return nil, fmt.Errorf("input string is empty")
    }

    var payload DataPayload
    err := json.Unmarshal([]byte(input), &payload)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    if payload.ID == "" {
        return nil, fmt.Errorf("field 'id' is required")
    }
    if payload.Value < 0 {
        return nil, fmt.Errorf("field 'value' must be non-negative")
    }

    return &payload, nil
}

func main() {
    // Example usage
    jsonStr := `{"id": "test-123", "value": 42, "tags": ["go", "json"]}`
    result, err := ParseAndValidateJSON(jsonStr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Parsed payload: %+v\n", result)
}