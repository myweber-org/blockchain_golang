package main

import (
    "fmt"
    "strings"
)

type DataCleaner struct {
    seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        seen: make(map[string]bool),
    }
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
    unique := []string{}
    for _, item := range items {
        trimmed := strings.TrimSpace(item)
        if trimmed == "" {
            continue
        }
        if !dc.seen[trimmed] {
            dc.seen[trimmed] = true
            unique = append(unique, trimmed)
        }
    }
    return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    return len(parts[0]) > 0 && len(parts[1]) > 0
}

func main() {
    cleaner := NewDataCleaner()
    
    data := []string{"  apple  ", "banana", "apple", "", "cherry", "banana"}
    cleaned := cleaner.RemoveDuplicates(data)
    fmt.Println("Deduplicated:", cleaned)
    
    emails := []string{"test@example.com", "invalid-email", "user@domain"}
    for _, email := range emails {
        fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
    }
}