package data

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

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

func FormatTimestamp(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func CalculateChecksum(data []byte) uint32 {
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}
	return sum % 256
}

func FilterSlice[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
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

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64 = records[0].Value
	count := len(records)

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func ValidateRecords(records []DataRecord) []error {
	var errors []error
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			errors = append(errors, fmt.Errorf("record '%s' has invalid ID: %d", record.Name, record.ID))
		}

		if seenIDs[record.ID] {
			errors = append(errors, fmt.Errorf("duplicate ID found: %d", record.ID))
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			errors = append(errors, fmt.Errorf("record '%s' has negative value: %.2f", record.Name, record.Value))
		}
	}

	return errors
}package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Value   int    `json:"value"`
	Active  bool   `json:"active"`
}

func processCSVData(reader io.Reader) ([]Record, error) {
	csvReader := csv.NewReader(reader)
	records := []Record{}

	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		record := Record{
			ID:     getField(row, headerMap, "id"),
			Name:   getField(row, headerMap, "name"),
			Active: strings.ToLower(getField(row, headerMap, "active")) == "true",
		}

		fmt.Sscanf(getField(row, headerMap, "value"), "%d", &record.Value)
		records = append(records, record)
	}

	return records, nil
}

func getField(row []string, headerMap map[string]int, key string) string {
	if idx, exists := headerMap[key]; exists && idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

func generateJSONOutput(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonData), nil
}

func main() {
	csvData := `id,name,value,active
1,alpha,100,true
2,beta,200,false
3,gamma,300,true`

	reader := strings.NewReader(csvData)
	records, err := processCSVData(reader)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		os.Exit(1)
	}

	jsonOutput, err := generateJSONOutput(records)
	if err != nil {
		fmt.Printf("Error generating output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processed Data:")
	fmt.Println(jsonOutput)
}