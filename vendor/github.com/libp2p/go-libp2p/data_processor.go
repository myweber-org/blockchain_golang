
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Active  bool    `json:"active"`
	Tags    []string `json:"tags"`
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.FieldsPerRecord = -1

	var records []Record
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		if len(row) < 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])
		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			continue
		}

		active := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		var tags []string
		if len(row) > 4 {
			tagStr := strings.TrimSpace(row[4])
			if tagStr != "" {
				tags = strings.Split(tagStr, "|")
			}
		}

		record := Record{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
			Tags:   tags,
		}

		records = append(records, record)
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func filterRecords(records []Record, predicate func(Record) bool) []Record {
	var filtered []Record
	for _, record := range records {
		if predicate(record) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func calculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 {
			min = record.Value
			max = record.Value
		} else {
			if record.Value < min {
				min = record.Value
			}
			if record.Value > max {
				max = record.Value
			}
		}
	}

	average := sum / float64(count)
	return average, min, max
}

func processDataFile(inputFile, outputFile string) error {
	records, err := parseCSVFile(inputFile)
	if err != nil {
		return err
	}

	activeRecords := filterRecords(records, func(r Record) bool {
		return r.Active
	})

	jsonOutput, err := convertToJSON(activeRecords)
	if err != nil {
		return err
	}

	avg, min, max := calculateStats(activeRecords)
	stats := fmt.Sprintf("Statistics:\n  Records: %d\n  Average: %.2f\n  Min: %.2f\n  Max: %.2f\n",
		len(activeRecords), avg, min, max)

	output := stats + "\nProcessed Data:\n" + jsonOutput

	return os.WriteFile(outputFile, []byte(output), 0644)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.txt>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	err := processDataFile(inputFile, outputFile)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %s to %s\n", inputFile, outputFile)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

	return records, nil
}

func parseRow(row []string, lineNum int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}
	record.Active = active

	return record, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d (must be positive)", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Name == "" {
			return fmt.Errorf("empty name for record ID: %d", record.ID)
		}

		if record.Value < 0 {
			return fmt.Errorf("negative value for record ID: %d", record.ID)
		}
	}

	return nil
}

func calculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	var maxValue float64

	for _, record := range records {
		sum += record.Value
		if record.Value > maxValue {
			maxValue = record.Value
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, maxValue, activeCount
}

func processDataFile(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	avg, max, activeCount := calculateStatistics(records)

	fmt.Printf("Data processing complete\n")
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
	fmt.Printf("Active records: %d\n", activeCount)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if err := processDataFile(filename); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}