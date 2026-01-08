package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// DataPayload represents a simple incoming data structure.
type DataPayload struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ValidatePayload checks if the required fields are present and valid.
func ValidatePayload(payload DataPayload) error {
	if payload.ID <= 0 {
		return fmt.Errorf("invalid ID: must be positive integer")
	}
	if payload.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if payload.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}

// ProcessJSONData parses a JSON byte slice and validates the content.
func ProcessJSONData(rawData []byte) (*DataPayload, error) {
	var payload DataPayload
	if err := json.Unmarshal(rawData, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if err := ValidatePayload(payload); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &payload, nil
}

func main() {
	// Example JSON data
	jsonData := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`

	processed, err := ProcessJSONData([]byte(jsonData))
	if err != nil {
		log.Fatalf("Error processing data: %v", err)
	}

	fmt.Printf("Processed payload: %+v\n", processed)
}