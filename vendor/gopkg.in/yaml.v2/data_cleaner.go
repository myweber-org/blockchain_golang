
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Valid bool
}

func DeduplicateEmails(emails []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, email := range emails {
        if _, exists := seen[email]; !exists {
            seen[email] = struct{}{}
            result = append(result, email)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        if ValidateEmail(record.Email) && !emailSet[record.Email] {
            emailSet[record.Email] = true
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    records := []DataRecord{
        {1, "user@example.com", false},
        {2, "invalid-email", false},
        {3, "user@example.com", false},
        {4, "another@test.org", false},
    }
    
    cleaned := CleanData(records)
    fmt.Printf("Original: %d records\n", len(records))
    fmt.Printf("Cleaned: %d records\n", len(cleaned))
    
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if len(email) < 5 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func cleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	for _, record := range records {
		if validateEmail(record.Email) {
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return deduplicateRecords(cleaned)
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "USER@example.com", false},
		{4, "test@domain.com", false},
		{5, "user@example.com", false},
	}

	cleaned := cleanData(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID    int
	Name  string
	Email string
	Score float64
}

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	lineNum := 1
	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row %d: %w", lineNum, err)
		}

		cleanedRow, err := validateAndCleanRow(row)
		if err != nil {
			fmt.Printf("Skipping row %d: %v\n", lineNum, err)
			continue
		}

		if err := writer.Write(cleanedRow); err != nil {
			return fmt.Errorf("error writing row %d: %w", lineNum, err)
		}
	}

	return nil
}

func validateAndCleanRow(row []string) ([]string, error) {
	if len(row) != 4 {
		return nil, fmt.Errorf("expected 4 columns, got %d", len(row))
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil || id <= 0 {
		return nil, fmt.Errorf("invalid ID: %s", row[0])
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return nil, fmt.Errorf("empty name")
	}
	name = strings.Title(strings.ToLower(name))

	email := strings.TrimSpace(row[2])
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email: %s", email)
	}
	email = strings.ToLower(email)

	score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil || score < 0 || score > 100 {
		return nil, fmt.Errorf("invalid score: %s", row[3])
	}

	return []string{
		strconv.Itoa(id),
		name,
		email,
		strconv.FormatFloat(score, 'f', 2, 64),
	}, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned data. Output written to %s\n", outputFile)
}