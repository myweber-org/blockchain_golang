package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateJSON(rawData []byte) (*UserData, error) {
	var user UserData
	err := json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonInput := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	user, err := ValidateJSON([]byte(jsonInput))
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}
	fmt.Printf("Valid user: %+v\n", user)
}
package data

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type Record struct {
	ID        string
	Email     string
	Timestamp time.Time
	Metadata  map[string]interface{}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func TransformRecord(record Record) (Record, error) {
	if err := ValidateEmail(record.Email); err != nil {
		return Record{}, err
	}

	normalizedEmail := NormalizeString(record.Email)
	processedMetadata := make(map[string]interface{})

	for key, value := range record.Metadata {
		switch v := value.(type) {
		case string:
			processedMetadata[key] = NormalizeString(v)
		default:
			processedMetadata[key] = v
		}
	}

	return Record{
		ID:        record.ID,
		Email:     normalizedEmail,
		Timestamp: record.Timestamp.UTC(),
		Metadata:  processedMetadata,
	}, nil
}

func FilterRecords(records []Record, predicate func(Record) bool) []Record {
	var filtered []Record
	for _, record := range records {
		if predicate(record) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func MergeMetadata(records []Record) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, record := range records {
		for key, value := range record.Metadata {
			if existing, exists := merged[key]; exists {
				if slice, ok := existing.([]interface{}); ok {
					merged[key] = append(slice, value)
				} else {
					merged[key] = []interface{}{existing, value}
				}
			} else {
				merged[key] = value
			}
		}
	}
	return merged
}