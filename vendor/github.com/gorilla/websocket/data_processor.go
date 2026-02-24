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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row length at line %d", lineNumber)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %v", lineNumber, err)
		}

		value, err := strconv.Atoi(row[2])
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %v", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return records, nil
}

func serializeToJSON(records []Record) (string, error) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func filterRecords(records []Record, minValue int) []Record {
	var filtered []Record
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func calculateStatistics(records []Record) (int, int, float64) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	total := 0
	maxValue := records[0].Value
	minValue := records[0].Value

	for _, record := range records {
		total += record.Value
		if record.Value > maxValue {
			maxValue = record.Value
		}
		if record.Value < minValue {
			minValue = record.Value
		}
	}

	average := float64(total) / float64(len(records))
	return maxValue, minValue, average
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

	fmt.Printf("Processed %d records\n", len(records))

	filtered := filterRecords(records, 50)
	fmt.Printf("Filtered to %d records with value >= 50\n", len(filtered))

	maxVal, minVal, avgVal := calculateStatistics(filtered)
	fmt.Printf("Statistics - Max: %d, Min: %d, Average: %.2f\n", maxVal, minVal, avgVal)

	jsonOutput, err := serializeToJSON(filtered)
	if err != nil {
		fmt.Printf("Error serializing to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nJSON Output:")
	fmt.Println(jsonOutput)
}