
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// ValidateJSONString checks if the provided string is valid JSON.
func ValidateJSONString(input string) (bool, error) {
    var js interface{}
    decoder := json.NewDecoder(strings.NewReader(input))
    decoder.DisallowUnknownFields()
    err := decoder.Decode(&js)
    if err != nil {
        return false, fmt.Errorf("invalid JSON: %w", err)
    }
    return true, nil
}

// ExtractField attempts to extract a string field from a JSON string by key.
func ExtractField(jsonStr, key string) (string, error) {
    var data map[string]interface{}
    err := json.Unmarshal([]byte(jsonStr), &data)
    if err != nil {
        return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    value, exists := data[key]
    if !exists {
        return "", fmt.Errorf("key '%s' not found in JSON object", key)
    }

    strValue, ok := value.(string)
    if !ok {
        return "", fmt.Errorf("value for key '%s' is not a string", key)
    }
    return strValue, nil
}

func main() {
    testJSON := `{"name": "Alice", "age": 30, "city": "London"}`

    valid, err := ValidateJSONString(testJSON)
    if valid {
        fmt.Println("JSON is valid.")
    } else {
        fmt.Printf("Validation error: %v\n", err)
    }

    name, err := ExtractField(testJSON, "name")
    if err != nil {
        fmt.Printf("Error extracting field: %v\n", err)
    } else {
        fmt.Printf("Extracted name: %s\n", name)
    }
}
package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}