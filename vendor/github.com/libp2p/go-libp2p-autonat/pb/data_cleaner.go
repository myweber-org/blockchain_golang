
package main

import (
	"fmt"
)

// RemoveDuplicates removes duplicate strings from a slice.
// It preserves the order of the first occurrence of each unique element.
func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	encountered := map[string]bool{}
	result := []DataRecord{}

	for _, record := range records {
		key := fmt.Sprintf("%d|%s|%s", record.ID, record.Name, record.Email)
		if !encountered[key] {
			encountered[key] = true
			result = append(result, record)
		}
	}
	return result
}

func ValidateRecords(records []DataRecord) []DataRecord {
	validated := []DataRecord{}
	for _, record := range records {
		record.Valid = record.ID > 0 &&
			len(strings.TrimSpace(record.Name)) > 0 &&
			strings.Contains(record.Email, "@")
		validated = append(validated, record)
	}
	return validated
}

func CleanData(records []DataRecord) []DataRecord {
	unique := RemoveDuplicates(records)
	validated := ValidateRecords(unique)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{1, "John Doe", "john@example.com", false},
		{3, "", "invalid-email", false},
	}

	cleaned := CleanData(sampleData)

	fmt.Println("Original records:", len(sampleData))
	fmt.Println("Cleaned records:", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Valid: %v\n", record.ID, record.Name, record.Valid)
	}
}