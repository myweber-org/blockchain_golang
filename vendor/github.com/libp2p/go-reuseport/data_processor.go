
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
	Delimiter  rune
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
		Delimiter:  ',',
	}
}

func (dp *DataProcessor) ValidateCSV() error {
	file, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = dp.Delimiter

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if len(headers) == 0 {
		return fmt.Errorf("csv file contains no headers")
	}

	rowCount := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row %d: %w", rowCount, err)
		}

		if len(record) != len(headers) {
			return fmt.Errorf("row %d has %d columns, expected %d", rowCount, len(record), len(headers))
		}

		for i, field := range record {
			if strings.TrimSpace(field) == "" {
				return fmt.Errorf("empty value in row %d, column %d", rowCount, i+1)
			}
		}
		rowCount++
	}

	return nil
}

func (dp *DataProcessor) TransformData() error {
	if err := dp.ValidateCSV(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.Comma = dp.Delimiter
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	transformedHeaders := make([]string, len(headers))
	for i, header := range headers {
		transformedHeaders[i] = strings.ToUpper(strings.TrimSpace(header))
	}

	if err := writer.Write(transformedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	rowCount := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row %d: %w", rowCount, err)
		}

		transformedRecord := make([]string, len(record))
		for i, field := range record {
			transformedRecord[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(transformedRecord); err != nil {
			return fmt.Errorf("failed to write row %d: %w", rowCount, err)
		}
		rowCount++
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	
	fmt.Printf("Processing %s to %s\n", processor.InputPath, processor.OutputPath)
	
	if err := processor.TransformData(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Data processing completed successfully")
}