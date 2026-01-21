
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
    if len(trimmed) < 3 || len(trimmed) > 254 {
        return false
    }
    atIndex := strings.LastIndex(trimmed, "@")
    if atIndex < 1 || atIndex == len(trimmed)-1 {
        return false
    }
    dotIndex := strings.LastIndex(trimmed[atIndex:], ".")
    if dotIndex < 1 || dotIndex == len(trimmed[atIndex:])-1 {
        return false
    }
    return true
}

func ProcessUserInput(username, email string) (*UserData, error) {
    normalizedUsername := NormalizeUsername(username)
    if normalizedUsername == "" {
        return nil, fmt.Errorf("invalid username: contains no valid characters")
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
    testData := []struct {
        username string
        email    string
    }{
        {"  John_Doe-123  ", "john@example.com"},
        {"Alice.Bob", "alice@test.org"},
        {"   ", "invalid-email"},
        {"Test-User_1", "bad@email"},
    }

    for _, td := range testData {
        user, err := ProcessUserInput(td.username, td.email)
        if err != nil {
            fmt.Printf("Error processing %s, %s: %v\n", td.username, td.email, err)
        } else {
            fmt.Printf("Processed: %+v\n", user)
        }
    }
}