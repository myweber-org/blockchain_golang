
package main

import (
    "encoding/csv"
    "errors"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
    Valid bool
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    // Skip header
    _, err = reader.Read()
    if err != nil {
        return nil, err
    }

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) < 4 {
            continue
        }

        id, idErr := strconv.Atoi(strings.TrimSpace(row[0]))
        name := strings.TrimSpace(row[1])
        value, valErr := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
        valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

        if idErr != nil || valErr != nil || name == "" {
            continue
        }

        record := DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
            Valid: valid,
        }
        records = append(records, record)
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
    if len(records) == 0 {
        return nil, errors.New("no records to validate")
    }

    validRecords := make([]DataRecord, 0)
    for _, record := range records {
        if record.Valid && record.Value >= 0 {
            validRecords = append(validRecords, record)
        }
    }

    return validRecords, nil
}

func CalculateAverage(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0.0
    }

    total := 0.0
    for _, record := range records {
        total += record.Value
    }

    return total / float64(len(records))
}
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

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
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

		if lineNumber == 1 {
			continue
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" || record.Email == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		if record.Active != "true" && record.Active != "false" {
			return nil, fmt.Errorf("invalid active status at line %d: %s", lineNumber, record.Active)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateEmailFormat(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if record.Active == "true" {
			active = append(active, record)
		}
	}
	return active
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	activeRecords := FilterActiveRecords(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Active records: %d\n", len(activeRecords))

	for i, record := range activeRecords {
		fmt.Printf("%d. %s <%s>\n", i+1, record.Name, record.Email)
	}
}