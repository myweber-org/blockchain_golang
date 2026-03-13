package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

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

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
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

	fmt.Printf("Successfully cleaned CSV: %s -> %s\n", inputFile, outputFile)
}
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

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		
		if validateEmail(cleanEmail) && !emailSet[cleanEmail] {
			emailSet[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	emails := []string{
		"test@example.com",
		"TEST@example.com",
		"invalid-email",
		"another@test.org",
		"test@example.com",
		"  spaced@domain.com  ",
	}

	fmt.Println("Original emails:", emails)
	deduped := deduplicateEmails(emails)
	fmt.Println("Deduplicated:", deduped)

	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "DUPLICATE@example.com", false},
		{3, "user@example.com", false},
		{4, "bad-format", false},
		{5, "new@test.co.uk", false},
	}

	cleanedRecords := processRecords(records)
	fmt.Printf("\nCleaned records: %d out of %d\n", len(cleanedRecords), len(records))
	
	for _, r := range cleanedRecords {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}