
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ProcessCSVFile(inputPath string, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	csvReader := csv.NewReader(inputFile)
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	processedHeaders := append(headers, "Processed", "Category")
	if err := csvWriter.Write(processedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) < 4 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			continue
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		active, err := strconv.ParseBool(row[3])
		if err != nil {
			continue
		}

		record := DataRecord{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		}

		processedValue := record.Value * 1.1
		category := determineCategory(record.Value)

		outputRow := []string{
			strconv.Itoa(record.ID),
			record.Name,
			strconv.FormatFloat(record.Value, 'f', 2, 64),
			strconv.FormatBool(record.Active),
			strconv.FormatFloat(processedValue, 'f', 2, 64),
			category,
		}

		if err := csvWriter.Write(outputRow); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}

		recordCount++
	}

	fmt.Printf("Processed %d records successfully\n", recordCount)
	return nil
}

func determineCategory(value float64) string {
	switch {
	case value < 10:
		return "Low"
	case value < 50:
		return "Medium"
	case value < 100:
		return "High"
	default:
		return "Premium"
	}
}

func ValidateCSVFormat(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return false, fmt.Errorf("failed to read headers: %w", err)
	}

	expectedHeaders := []string{"ID", "Name", "Value", "Active"}
	if len(headers) != len(expectedHeaders) {
		return false, nil
	}

	for i, header := range headers {
		if strings.ToLower(header) != strings.ToLower(expectedHeaders[i]) {
			return false, nil
		}
	}

	return true, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	valid, err := ValidateCSVFormat(inputFile)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	if !valid {
		fmt.Println("Invalid CSV format")
		os.Exit(1)
	}

	if err := ProcessCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data processing completed successfully")
}