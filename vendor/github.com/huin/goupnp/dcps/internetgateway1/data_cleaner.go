
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

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return len(email) > 5
}

func validateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = validateEmail(record.Email)
		validated = append(validated, record)
	}
	return validated
}

func printRecords(records []DataRecord) {
	fmt.Println("Processed Records:")
	for _, record := range records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Status: %s\n",
			record.ID, record.Name, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.org", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
		{5, "Alice Brown", "alice@test.co", false},
	}

	fmt.Println("Original records:", len(records))
	unique := deduplicateRecords(records)
	fmt.Println("After deduplication:", len(unique))

	validated := validateRecords(unique)
	printRecords(validated)
}
package main

import "fmt"

func removeDuplicates(input []int) []int {
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
	cleaned := removeDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    if len(nums) == 0 {
        return nums
    }

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
    data := []int{1, 2, 3, 2, 4, 5, 1, 6, 3}
    cleaned := RemoveDuplicates(data)
    fmt.Println("Original:", data)
    fmt.Println("Cleaned:", cleaned)
}