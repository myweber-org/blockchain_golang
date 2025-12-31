package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func ValidateJSON(data []byte) (bool, error) {
	if !json.Valid(data) {
		return false, fmt.Errorf("invalid JSON structure")
	}
	return true, nil
}

func ParseUserData(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	jsonData := []byte(`{"name": "Alice", "age": 30, "active": true}`)

	valid, err := ValidateJSON(jsonData)
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}
	fmt.Printf("JSON is valid: %v\n", valid)

	parsed, err := ParseUserData(jsonData)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	fmt.Printf("Parsed data: %v\n", parsed)
}