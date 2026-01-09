
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Tags      []string  `json:"tags"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func TransformTags(tags []string) []string {
	transformed := make([]string, 0, len(tags))
	for _, tag := range tags {
		cleanTag := strings.TrimSpace(tag)
		if cleanTag != "" {
			transformed = append(transformed, strings.ToLower(cleanTag))
		}
	}
	return transformed
}

func ProcessRecord(record DataRecord) (DataRecord, error) {
	if !ValidateEmail(record.Email) {
		return DataRecord{}, fmt.Errorf("invalid email format: %s", record.Email)
	}

	if record.Value < 0 {
		record.Value = 0
	}

	record.Tags = TransformTags(record.Tags)
	record.Timestamp = time.Now().UTC()

	return record, nil
}

func SerializeRecord(record DataRecord) ([]byte, error) {
	return json.MarshalIndent(record, "", "  ")
}

func main() {
	sampleRecord := DataRecord{
		ID:        "user-123",
		Email:     "test@example.com",
		Timestamp: time.Now(),
		Value:     42.5,
		Tags:      []string{"  Go  ", "BACKEND", "", "Data  "},
	}

	processed, err := ProcessRecord(sampleRecord)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	jsonData, err := SerializeRecord(processed)
	if err != nil {
		fmt.Printf("Serialization error: %v\n", err)
		return
	}

	fmt.Printf("Processed record:\n%s\n", string(jsonData))
}