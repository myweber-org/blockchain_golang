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
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
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

		if lineNumber == 1 {
			continue
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		if !strings.Contains(record.Email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	emailSet := make(map[string]bool)

	for _, record := range records {
		if emailSet[record.Email] {
			return fmt.Errorf("duplicate email found: %s", record.Email)
		}
		emailSet[record.Email] = true

		if record.Active != "true" && record.Active != "false" {
			return fmt.Errorf("invalid active status for record %s: %s", record.ID, record.Active)
		}
	}

	return nil
}

func GenerateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}