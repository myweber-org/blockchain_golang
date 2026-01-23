
package main

import (
    "fmt"
    "strings"
)

func RemoveDuplicates(data []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, item := range data {
        if !seen[item] {
            seen[item] = true
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
    if len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    return true
}

func CleanData(data []string) []string {
    uniqueData := RemoveDuplicates(data)
    cleaned := []string{}
    for _, item := range uniqueData {
        if ValidateEmail(item) {
            cleaned = append(cleaned, item)
        }
    }
    return cleaned
}

func main() {
    sampleData := []string{
        "user@example.com",
        "admin@test.org",
        "user@example.com",
        "invalid-email",
        "another@domain.com",
        "",
        "missing@domain",
        "another@domain.com",
    }

    fmt.Println("Original data:", sampleData)
    cleaned := CleanData(sampleData)
    fmt.Println("Cleaned data:", cleaned)
}