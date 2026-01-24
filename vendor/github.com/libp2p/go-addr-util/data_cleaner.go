package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Phone string
}

func generateHash(email, phone string) string {
	data := strings.ToLower(strings.TrimSpace(email)) + "|" + strings.TrimSpace(phone)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func deduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		hash := generateHash(record.Email, record.Phone)
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []Record) []Record {
	var valid []Record
	for _, record := range records {
		if validateEmail(record.Email) && len(record.Phone) >= 10 {
			valid = append(valid, record)
		}
	}
	return deduplicateRecords(valid)
}

func main() {
	sampleData := []Record{
		{ID: "1", Email: "user@example.com", Phone: "1234567890"},
		{ID: "2", Email: "user@example.com", Phone: "1234567890"},
		{ID: "3", Email: "invalid-email", Phone: "1234567890"},
		{ID: "4", Email: "another@test.org", Phone: "0987654321"},
		{ID: "5", Email: "another@test.org", Phone: "0987654321"},
		{ID: "6", Email: "short@phone.com", Phone: "123"},
	}

	cleaned := cleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))

	for _, record := range cleaned {
		fmt.Printf("ID: %s, Email: %s\n", record.ID, record.Email)
	}
}package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Cleaned:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Cleaned:", uniqueStrings)
}