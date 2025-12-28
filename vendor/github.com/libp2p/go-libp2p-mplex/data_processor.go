package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Active  string
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) < 4 {
			continue
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" {
			continue
		}

		if !validateEmail(record.Email) {
			fmt.Printf("Warning: Invalid email format for record ID %s\n", record.ID)
		}

		records = append(records, record)
	}

	return records, nil
}

func filterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if strings.ToLower(record.Active) == "true" || record.Active == "1" {
			active = append(active, record)
		}
	}
	return active
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	activeRecords := filterActiveRecords(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Active records: %d\n", len(activeRecords))

	for i, record := range activeRecords {
		if i < 5 {
			fmt.Printf("Active user: %s (%s)\n", record.Name, record.Email)
		}
	}
}