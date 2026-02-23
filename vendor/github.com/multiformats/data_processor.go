
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?-]+$`)
	return &DataProcessor{allowedPattern: pattern}
}

func (dp *DataProcessor) SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	if !dp.allowedPattern.MatchString(trimmed) {
		return "", false
	}

	return trimmed, true
}

func (dp *DataProcessor) ProcessUserData(rawData string) (string, error) {
	sanitized, valid := dp.SanitizeInput(rawData)
	if !valid {
		return "", &InvalidInputError{Input: rawData}
	}

	processed := strings.ToUpper(sanitized)
	return processed, nil
}

type InvalidInputError struct {
	Input string
}

func (e *InvalidInputError) Error() string {
	return "input contains invalid characters or is empty"
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessCSVFile(inputPath string) ([]DataRecord, error) {
	file, err := os.Open(inputPath)
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

		if len(row) < 4 {
			continue
		}

		record, err := parseRow(row)
		if err != nil {
			fmt.Printf("Warning: skipping invalid row at line %d: %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID format: %w", err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value format: %w", err)
	}
	record.Value = value

	validStr := strings.ToLower(strings.TrimSpace(row[3]))
	record.Valid = validStr == "true" || validStr == "1"

	return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total records processed: %d\n", len(records))

	validRecords := FilterValidRecords(records)
	fmt.Printf("Valid records: %d\n", len(validRecords))

	average := CalculateAverage(validRecords)
	fmt.Printf("Average value of valid records: %.2f\n", average)
}