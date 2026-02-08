
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSVProcessor struct {
	FilePath string
	Headers  []string
	Records  [][]string
}

func NewCSVProcessor(filePath string) *CSVProcessor {
	return &CSVProcessor{
		FilePath: filePath,
	}
}

func (cp *CSVProcessor) Read() error {
	file, err := os.Open(cp.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}
	cp.Headers = headers

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}
		cp.Records = append(cp.Records, record)
	}
	return nil
}

func (cp *CSVProcessor) Write(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(cp.Headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, record := range cp.Records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}
	return nil
}

func (cp *CSVProcessor) FilterByColumn(columnIndex int, filterFunc func(string) bool) [][]string {
	var filtered [][]string
	for _, record := range cp.Records {
		if columnIndex < len(record) && filterFunc(record[columnIndex]) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	processor := NewCSVProcessor("input.csv")
	if err := processor.Read(); err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	fmt.Printf("Read %d records with headers: %v\n", len(processor.Records), processor.Headers)

	filtered := processor.FilterByColumn(0, func(value string) bool {
		return len(value) > 0 && value != "exclude"
	})
	fmt.Printf("Filtered to %d records\n", len(filtered))

	if err := processor.Write("output.csv"); err != nil {
		fmt.Printf("Error writing CSV: %v\n", err)
	}
}