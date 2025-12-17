package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON data into a generic map.
func ParseUserData(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := `{"name": "alice", "age": 30, "active": true}`

	// Validate the JSON
	isValid, err := ValidateJSON([]byte(sampleJSON))
	if err != nil {
		log.Printf("Validation error: %v", err)
	} else {
		fmt.Printf("JSON is valid: %t\n", isValid)
	}

	// Parse the JSON
	parsedData, err := ParseUserData(sampleJSON)
	if err != nil {
		log.Printf("Parse error: %v", err)
	} else {
		fmt.Printf("Parsed data: %+v\n", parsedData)
	}
}