
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    return result
}

func main() {
    data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
    cleaned := RemoveDuplicates(data)
    fmt.Println("Original:", data)
    fmt.Println("Cleaned:", cleaned)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, rec := range records {
		email := strings.ToLower(strings.TrimSpace(rec.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, rec)
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func CleanData(records []Record) ([]Record, error) {
	var cleaned []Record
	for _, rec := range records {
		if err := ValidateEmail(rec.Email); err != nil {
			continue
		}
		cleaned = append(cleaned, rec)
	}
	return DeduplicateRecords(cleaned), nil
}

func main() {
	sampleData := []Record{
		{1, "user@example.com", true},
		{2, "USER@example.com", true},
		{3, "test@domain.org", true},
		{4, "invalid-email", true},
		{5, "", true},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", rec.ID, rec.Email)
	}
}