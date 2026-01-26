
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
	ID        int
	Name      string
	Email     string
	Age       int
	Active    bool
	Score     float64
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	var data []DataRecord
	for i, row := range records {
		if i == 0 {
			continue
		}

		if len(row) < 6 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])
		email := strings.TrimSpace(row[2])

		age, err := strconv.Atoi(strings.TrimSpace(row[3]))
		if err != nil {
			continue
		}

		active := strings.ToLower(strings.TrimSpace(row[4])) == "true"

		score, err := strconv.ParseFloat(strings.TrimSpace(row[5]), 64)
		if err != nil {
			continue
		}

		record := DataRecord{
			ID:     id,
			Name:   name,
			Email:  email,
			Age:    age,
			Active: active,
			Score:  score,
		}

		if validateRecord(record) {
			data = append(data, record)
		}
	}

	return data, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID <= 0 {
		return false
	}

	if record.Name == "" || len(record.Name) > 100 {
		return false
	}

	if !strings.Contains(record.Email, "@") {
		return false
	}

	if record.Age < 0 || record.Age > 150 {
		return false
	}

	if record.Score < 0 || record.Score > 100 {
		return false
	}

	return true
}

func cleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if seenIDs[record.ID] {
			continue
		}

		record.Name = strings.Title(strings.ToLower(record.Name))
		record.Email = strings.ToLower(record.Email)

		if record.Score < 50 {
			record.Active = false
		}

		cleaned = append(cleaned, record)
		seenIDs[record.ID] = true
	}

	return cleaned
}

func writeCSVFile(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Email", "Age", "Active", "Score"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, record := range records {
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			record.Email,
			strconv.Itoa(record.Age),
			strconv.FormatBool(record.Active),
			strconv.FormatFloat(record.Score, 'f', 2, 64),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

func processDataFile(inputFile, outputFile string) error {
	records, err := parseCSVFile(inputFile)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	fmt.Printf("Parsed %d records\n", len(records))

	cleaned := cleanData(records)
	fmt.Printf("Cleaned to %d records\n", len(cleaned))

	if err := writeCSVFile(outputFile, cleaned); err != nil {
		return fmt.Errorf("writing failed: %w", err)
	}

	fmt.Printf("Data written to %s\n", outputFile)
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processDataFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}