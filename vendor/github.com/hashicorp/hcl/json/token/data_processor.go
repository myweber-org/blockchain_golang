
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row length at line %d", lineNumber)
		}

		age, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("invalid age at line %d: %w", lineNumber, err)
		}

		score, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid score at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			Name:  row[0],
			Age:   age,
			Score: score,
		})
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshaling error: %w", err)
	}
	return string(jsonData), nil
}

func calculateAverageScore(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0.0
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	jsonOutput, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processed Records:")
	fmt.Println(jsonOutput)
	fmt.Printf("\nTotal Records: %d\n", len(records))
	fmt.Printf("Average Score: %.2f\n", calculateAverageScore(records))
}package data

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Email     string
	Timestamp time.Time
	Value     float64
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

func ProcessRecord(record DataRecord) (DataRecord, error) {
	if err := ValidateEmail(record.Email); err != nil {
		return DataRecord{}, err
	}

	record.Email = NormalizeString(record.Email)
	record.ID = strings.ToUpper(record.ID)

	if record.Value < 0 {
		record.Value = 0
	}

	return record, nil
}

func CalculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}