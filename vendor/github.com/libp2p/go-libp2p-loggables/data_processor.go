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
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	avg, max, count := CalculateStatistics(records)
	fmt.Printf("Processed %d records\n", count)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}