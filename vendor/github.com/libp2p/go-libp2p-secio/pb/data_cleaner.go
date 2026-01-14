
package main

import (
    "fmt"
    "strings"
)

func DeduplicateStrings(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
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
    if !strings.Contains(parts[1], ".") {
        return false
    }
    return true
}

func main() {
    emails := []string{
        "test@example.com",
        "user@domain.org",
        "test@example.com",
        "invalid-email",
        "another@test.net",
        "user@domain.org",
    }

    uniqueEmails := DeduplicateStrings(emails)
    fmt.Println("Unique emails:", uniqueEmails)

    for _, email := range uniqueEmails {
        if ValidateEmail(email) {
            fmt.Printf("%s is valid\n", email)
        } else {
            fmt.Printf("%s is invalid\n", email)
        }
    }
}package main

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
	cleanedData := removeDuplicates(data)
	fmt.Println("Original data:", data)
	fmt.Println("Cleaned data:", cleanedData)
}