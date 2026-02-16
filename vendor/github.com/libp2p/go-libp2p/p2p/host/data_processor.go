
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}
	if len(record.Tags) == 0 {
		return errors.New("record must have at least one tag")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Timestamp: record.Timestamp.UTC(),
		Value:     record.Value * multiplier,
		Tags:      append([]string{"processed"}, record.Tags...),
	}
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record, 1.5))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Timestamp: time.Now(),
			Value:     100.0,
			Tags:      []string{"test", "sample"},
		},
		{
			ID:        "rec002",
			Timestamp: time.Now().Add(-time.Hour),
			Value:     200.0,
			Tags:      []string{"production"},
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, record := range processed {
		fmt.Printf("Processed: %+v\n", record)
	}
}
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
	Valid bool    `json:"valid"`
}

func processCSVFile(inputPath string) ([]Record, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	headerSkipped := false

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		if len(row) < 4 {
			continue
		}

		age, _ := strconv.Atoi(row[1])
		score, _ := strconv.ParseFloat(row[2], 64)
		valid, _ := strconv.ParseBool(row[3])

		record := Record{
			Name:  row[0],
			Age:   age,
			Score: score,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func convertToJSON(records []Record) ([]byte, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json marshaling error: %w", err)
	}
	return jsonData, nil
}

func writeOutputFile(data []byte, outputPath string) error {
	err := os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.json>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	records, err := processCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	err = writeOutputFile(jsonData, outputFile)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records to %s\n", len(records), outputFile)
}