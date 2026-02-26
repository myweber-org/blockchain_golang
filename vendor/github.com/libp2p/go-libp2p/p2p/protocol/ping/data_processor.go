
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID    string
	Name  string
	Email string
	Valid bool
}

func processCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
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

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(row[0], row[1], row[2]),
		}

		if record.Valid {
			records = append(records, record)
		}
	}

	return records, nil
}

func validateRecord(id, name, email string) bool {
	if id == "" || name == "" || email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}

func generateReport(records []DataRecord) {
	fmt.Printf("Total valid records: %d\n", len(records))
	for _, record := range records {
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file_path>")
		os.Exit(1)
	}

	records, err := processCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	generateReport(records)
}
package main

import (
	"regexp"
	"strings"
)

// CleanString removes extra whitespace and normalizes line endings
func CleanString(input string) string {
	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim leading/trailing whitespace
	cleaned = strings.TrimSpace(cleaned)
	
	// Normalize line endings to Unix style
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")
	
	return cleaned
}

// NormalizeWhitespace ensures consistent spacing around punctuation
func NormalizeWhitespace(input string) string {
	// Add space after punctuation if missing
	re := regexp.MustCompile(`([.,!?])([^\s])`)
	normalized := re.ReplaceAllString(input, "$1 $2")
	
	// Remove space before punctuation
	re = regexp.MustCompile(`\s+([.,!?])`)
	normalized = re.ReplaceAllString(normalized, "$1")
	
	return normalized
}

// ProcessInput applies all cleaning and normalization steps
func ProcessInput(input string) string {
	cleaned := CleanString(input)
	normalized := NormalizeWhitespace(cleaned)
	return normalized
}