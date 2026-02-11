package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func TransformToTitleCase(input string) string {
	if input == "" {
		return input
	}
	return strings.ToUpper(input[:1]) + strings.ToLower(input[1:])
}

func PrettyPrintJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	email := "test@example.com"
	fmt.Printf("Email %s valid: %v\n", email, ValidateEmail(email))

	name := "john doe"
	fmt.Printf("Title case of '%s': %s\n", name, TransformToTitleCase(name))

	sample := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
	}
	pretty, _ := PrettyPrintJSON(sample)
	fmt.Println("Pretty JSON:")
	fmt.Println(pretty)
}