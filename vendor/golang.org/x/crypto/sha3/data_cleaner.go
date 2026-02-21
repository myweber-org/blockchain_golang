
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
	return true
}

func processRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord
	
	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		
		if !emailSet[cleanEmail] && validateEmail(cleanEmail) {
			emailSet[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	emails := []string{"test@example.com", "TEST@example.com", "invalid", "another@test.org"}
	unique := deduplicateEmails(emails)
	fmt.Printf("Unique emails: %v\n", unique)
	
	records := []DataRecord{
		{1, "user@domain.com", false},
		{2, "USER@domain.com", false},
		{3, "bad-email", false},
	}
	
	cleaned := processRecords(records)
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
}package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, line := range input {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	lines, err := readLines(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	cleaned := cleanData(lines)

	err = writeLines(outputFile, cleaned)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cleaned %d lines to %d unique lines\n", len(lines), len(cleaned))
}