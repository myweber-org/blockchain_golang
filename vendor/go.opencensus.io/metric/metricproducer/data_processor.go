
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func TrimAndTitle(s string) string {
	return strings.Title(strings.TrimSpace(s))
}

func ToJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	email := "test@example.com"
	fmt.Printf("Email %s valid: %v\n", email, ValidateEmail(email))

	name := "  john doe  "
	fmt.Printf("Original: '%s', Processed: '%s'\n", name, TrimAndTitle(name))

	person := map[string]string{
		"name":  "Alice",
		"email": "alice@example.com",
	}
	jsonStr, _ := ToJSON(person)
	fmt.Println("JSON output:", jsonStr)
}