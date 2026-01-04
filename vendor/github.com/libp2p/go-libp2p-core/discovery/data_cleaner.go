
package main

import (
    "fmt"
    "strings"
)

// DataCleaner provides methods for cleaning string data
type DataCleaner struct{}

// Deduplicate removes duplicate entries from a slice of strings
func (dc *DataCleaner) Deduplicate(items []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range items {
        trimmed := strings.TrimSpace(item)
        if trimmed != "" && !seen[trimmed] {
            seen[trimmed] = true
            result = append(result, trimmed)
        }
    }
    return result
}

// ValidateEmail checks if a string is a valid email format
func (dc *DataCleaner) ValidateEmail(email string) bool {
    if len(email) < 3 || len(email) > 254 {
        return false
    }
    
    atIndex := strings.Index(email, "@")
    if atIndex < 1 || atIndex == len(email)-1 {
        return false
    }
    
    dotIndex := strings.LastIndex(email[atIndex:], ".")
    if dotIndex < 1 || dotIndex == len(email[atIndex:])-1 {
        return false
    }
    
    return true
}

// NormalizeWhitespace replaces multiple spaces with single space
func (dc *DataCleaner) NormalizeWhitespace(text string) string {
    return strings.Join(strings.Fields(text), " ")
}

func main() {
    cleaner := &DataCleaner{}
    
    // Example usage
    data := []string{"  apple", "banana", "apple", "  ", "banana", "cherry"}
    cleaned := cleaner.Deduplicate(data)
    fmt.Printf("Deduplicated: %v\n", cleaned)
    
    email := "test@example.com"
    fmt.Printf("Email valid: %v\n", cleaner.ValidateEmail(email))
    
    text := "This   has   extra   spaces"
    fmt.Printf("Normalized: %s\n", cleaner.NormalizeWhitespace(text))
}