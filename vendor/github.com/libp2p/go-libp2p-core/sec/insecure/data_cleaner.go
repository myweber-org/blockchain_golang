
package main

import (
    "fmt"
    "strings"
)

func DeduplicateStrings(slice []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, item := range slice {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    if parts[0] == "" || parts[1] == "" {
        return false
    }
    if !strings.Contains(parts[1], ".") {
        return false
    }
    return true
}

func CleanPhoneNumber(phone string) string {
    var builder strings.Builder
    for _, ch := range phone {
        if ch >= '0' && ch <= '9' {
            builder.WriteRune(ch)
        }
    }
    return builder.String()
}

func main() {
    emails := []string{"test@example.com", "invalid-email", "test@example.com", "another@domain.org"}
    fmt.Println("Original emails:", emails)
    uniqueEmails := DeduplicateStrings(emails)
    fmt.Println("Deduplicated emails:", uniqueEmails)

    for _, email := range uniqueEmails {
        fmt.Printf("Email %s is valid: %v\n", email, ValidateEmail(email))
    }

    phone := "+1 (234) 567-8900"
    cleaned := CleanPhoneNumber(phone)
    fmt.Printf("Original phone: %s -> Cleaned: %s\n", phone, cleaned)
}