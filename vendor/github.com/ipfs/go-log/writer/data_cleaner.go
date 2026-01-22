package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataCleaner struct {
	InputPath  string
	OutputPath string
	Delimiter  rune
}

func NewDataCleaner(input, output string) *DataCleaner {
	return &DataCleaner{
		InputPath:  input,
		OutputPath: output,
		Delimiter:  ',',
	}
}

func (dc *DataCleaner) Clean() error {
	inputFile, err := os.Open(dc.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dc.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.Comma = dc.Delimiter
	writer := csv.NewWriter(outputFile)
	writer.Comma = dc.Delimiter

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := dc.cleanRow(headers)
	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := dc.cleanRow(record)
		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
		recordCount++
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	fmt.Printf("Cleaned %d records from %s to %s\n", recordCount, dc.InputPath, dc.OutputPath)
	return nil
}

func (dc *DataCleaner) cleanRow(row []string) []string {
	cleaned := make([]string, len(row))
	for i, value := range row {
		cleaned[i] = strings.TrimSpace(value)
		cleaned[i] = strings.ToUpper(cleaned[i])
	}
	return cleaned
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	cleaner := NewDataCleaner(os.Args[1], os.Args[2])
	if err := cleaner.Clean(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}