package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

// ParseUserData attempts to parse JSON data into a map.
func ParseUserData(jsonData []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := []byte(`{"name": "Alice", "age": 30, "active": true}`)

	if ValidateJSON(sampleJSON) {
		fmt.Println("JSON is valid.")
		parsed, err := ParseUserData(sampleJSON)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Parsed data: %v\n", parsed)
	} else {
		fmt.Println("Invalid JSON.")
	}
}