package datautils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&;]`)
	sanitized := re.ReplaceAllString(trimmed, "")
	return sanitized
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ContainsOnlyAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}
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

type Record struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func ProcessCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.Comment = '#'
    reader.FieldsPerRecord = 4

    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("line %d: %w", lineNumber, err)
        }

        record, err := parseRow(row)
        if err != nil {
            return nil, fmt.Errorf("line %d: %w", lineNumber, err)
        }

        if err := validateRecord(record); err != nil {
            return nil, fmt.Errorf("line %d: %w", lineNumber, err)
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found")
    }

    return records, nil
}

func parseRow(row []string) (Record, error) {
    var record Record

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID: %v", err)
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, errors.New("name cannot be empty")
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid value: %v", err)
    }
    record.Value = value

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, fmt.Errorf("invalid active flag: %v", err)
    }
    record.Active = active

    return record, nil
}

func validateRecord(r Record) error {
    if r.ID <= 0 {
        return errors.New("ID must be positive")
    }
    if r.Value < 0 {
        return errors.New("value cannot be negative")
    }
    if len(r.Name) > 100 {
        return errors.New("name too long")
    }
    return nil
}

func CalculateStats(records []Record) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var activeCount int
    var minValue float64 = records[0].Value

    for _, record := range records {
        sum += record.Value
        if record.Value < minValue {
            minValue = record.Value
        }
        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, minValue, activeCount
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) < 3 {
			return nil, fmt.Errorf("invalid row length: %d", len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return records, nil
}

func validateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID: %d", record.ID)
		}
		if record.Name == "" {
			return fmt.Errorf("empty name for ID: %d", record.ID)
		}
		if record.Value < 0 {
			return fmt.Errorf("negative value for ID: %d", record.ID)
		}
		if seen[record.ID] {
			return fmt.Errorf("duplicate ID: %d", record.ID)
		}
		seen[record.ID] = true
	}
	return nil
}

func calculateTotalValue(records []Record) float64 {
	var total float64
	for _, record := range records {
		total += record.Value
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := processCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	if err := validateRecords(records); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	total := calculateTotalValue(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", total)
}