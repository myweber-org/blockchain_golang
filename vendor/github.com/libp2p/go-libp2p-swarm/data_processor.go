
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	FilePath string
	Headers  []string
	Records  [][]string
}

func NewDataProcessor(filePath string) *DataProcessor {
	return &DataProcessor{
		FilePath: filePath,
	}
}

func (dp *DataProcessor) LoadAndValidate() error {
	file, err := os.Open(dp.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}
	dp.Headers = headers

	dp.Records = make([][]string, 0)
	lineNum := 2
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading line %d: %w", lineNum, err)
		}

		if len(record) != len(headers) {
			return fmt.Errorf("column count mismatch on line %d", lineNum)
		}

		for i, value := range record {
			record[i] = strings.TrimSpace(value)
			if record[i] == "" {
				return fmt.Errorf("empty value in column '%s' line %d", headers[i], lineNum)
			}
		}

		dp.Records = append(dp.Records, record)
		lineNum++
	}

	if len(dp.Records) == 0 {
		return fmt.Errorf("no data records found")
	}

	return nil
}

func (dp *DataProcessor) GetColumnStats(columnIndex int) (min, max string, err error) {
	if columnIndex < 0 || columnIndex >= len(dp.Headers) {
		return "", "", fmt.Errorf("invalid column index")
	}

	if len(dp.Records) == 0 {
		return "", "", fmt.Errorf("no records available")
	}

	minVal := dp.Records[0][columnIndex]
	maxVal := dp.Records[0][columnIndex]

	for _, record := range dp.Records[1:] {
		val := record[columnIndex]
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	return minVal, maxVal, nil
}

func (dp *DataProcessor) FilterRecords(filterFunc func([]string) bool) [][]string {
	filtered := make([][]string, 0)
	for _, record := range dp.Records {
		if filterFunc(record) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1])
	if err := processor.LoadAndValidate(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully loaded %d records with %d columns\n", 
		len(processor.Records), len(processor.Headers))

	for i, header := range processor.Headers {
		min, max, err := processor.GetColumnStats(i)
		if err == nil {
			fmt.Printf("Column '%s': min='%s', max='%s'\n", header, min, max)
		}
	}
}