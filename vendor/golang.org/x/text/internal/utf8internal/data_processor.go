package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func ValidateJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}
	return result, nil
}

func main() {
	sampleJSON := `{"name": "test", "value": 42}`
	parsed, err := ValidateJSON([]byte(sampleJSON))
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}
	fmt.Printf("Parsed data: %v\n", parsed)
}