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

	if len(result) == 0 {
		return nil, fmt.Errorf("JSON object is empty")
	}

	return result, nil
}

func main() {
	jsonData := `{"name": "test", "value": 123}`
	parsed, err := ValidateJSON([]byte(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed data: %v\n", parsed)
}