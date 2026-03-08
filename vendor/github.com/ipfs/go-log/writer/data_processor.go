
package main

import (
    "encoding/csv"
    "errors"
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

func ParseCSVFile(filename string) ([]DataRecord, error) {
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

        if len(row) != 4 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
        }

        record, err := parseRow(row, lineNumber)
        if err != nil {
            return nil, err
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func parseRow(row []string, lineNum int) (DataRecord, error) {
    var record DataRecord

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID at line %d: %v", lineNum, err)
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, fmt.Errorf("empty name at line %d", lineNum)
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid value at line %d: %v", lineNum, err)
    }
    record.Value = value

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, fmt.Errorf("invalid active flag at line %d: %v", lineNum, err)
    }
    record.Active = active

    return record, nil
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
    var active []DataRecord
    for _, record := range records {
        if record.Active {
            active = append(active, record)
        }
    }
    return active
}

func CalculateAverageValue(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0.0
    }

    total := 0.0
    for _, record := range records {
        total += record.Value
    }
    return total / float64(len(records))
}

func ValidateRecords(records []DataRecord) error {
    seenIDs := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid record ID: %d", record.ID)
        }
        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value for record ID %d", record.ID)
        }
    }
    return nil
}package main

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

func (dp *DataProcessor) ValidateRow(row []string) bool {
	if len(row) == 0 {
		return false
	}
	for _, field := range row {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

func (dp *DataProcessor) CleanField(field string) string {
	cleaned := strings.TrimSpace(field)
	cleaned = strings.ToUpper(cleaned)
	return cleaned
}

func (dp *DataProcessor) Process() error {
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
	writer.Comma = dp.Delimiter
	defer writer.Flush()

	lineCount := 0
	validCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		lineCount++

		if !dp.ValidateRow(record) {
			continue
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = dp.CleanField(field)
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing CSV: %w", err)
		}

		validCount++
	}

	fmt.Printf("Processed %d lines, %d valid records written to %s\n", 
		lineCount, validCount, dp.OutputPath)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}