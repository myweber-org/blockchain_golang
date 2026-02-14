
package main

import (
    "fmt"
    "strings"
)

// DataCleaner provides methods for cleaning string data
type DataCleaner struct{}

// Deduplicate removes duplicate entries from a slice of strings
func (dc *DataCleaner) Deduplicate(items []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range items {
        trimmed := strings.TrimSpace(item)
        if trimmed != "" && !seen[trimmed] {
            seen[trimmed] = true
            result = append(result, trimmed)
        }
    }
    return result
}

// ValidateEmail checks if a string is a valid email format
func (dc *DataCleaner) ValidateEmail(email string) bool {
    if len(email) < 3 || len(email) > 254 {
        return false
    }
    
    atIndex := strings.Index(email, "@")
    if atIndex < 1 || atIndex == len(email)-1 {
        return false
    }
    
    dotIndex := strings.LastIndex(email[atIndex:], ".")
    if dotIndex < 1 || dotIndex == len(email[atIndex:])-1 {
        return false
    }
    
    return true
}

// NormalizeWhitespace replaces multiple spaces with single space
func (dc *DataCleaner) NormalizeWhitespace(text string) string {
    return strings.Join(strings.Fields(text), " ")
}

func main() {
    cleaner := &DataCleaner{}
    
    // Example usage
    data := []string{"  apple", "banana", "apple", "  ", "banana", "cherry"}
    cleaned := cleaner.Deduplicate(data)
    fmt.Printf("Deduplicated: %v\n", cleaned)
    
    email := "test@example.com"
    fmt.Printf("Email valid: %v\n", cleaner.ValidateEmail(email))
    
    text := "This   has   extra   spaces"
    fmt.Printf("Normalized: %s\n", cleaner.NormalizeWhitespace(text))
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple ", " BANANA", "banana", "Cherry"}
	fmt.Println("Original:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated:", deduped)
	
	cleaner.Reset()
	testValue := "Test"
	fmt.Printf("Is '%s' duplicate? %v\n", testValue, cleaner.IsDuplicate(testValue))
	fmt.Printf("Is '%s' duplicate again? %v\n", testValue, cleaner.IsDuplicate(testValue))
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func cleanCSVData(inputPath, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read record: %w", err)
        }

        cleaned := make([]string, len(record))
        for i, field := range record {
            cleaned[i] = strings.TrimSpace(field)
            if cleaned[i] == "" {
                cleaned[i] = "N/A"
            }
        }

        if err := writer.Write(cleaned); err != nil {
            return fmt.Errorf("failed to write cleaned record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    if err := cleanCSVData(os.Args[1], os.Args[2]); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
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

func DeduplicateEmails(emails []string) []string {
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

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func CleanRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	cleaned := []DataRecord{}
	for _, record := range records {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		if ValidateEmail(record.Email) && !emailSet[record.Email] {
			record.Valid = true
			emailSet[record.Email] = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "test@example.com", false},
		{2, "invalid-email", false},
		{3, "TEST@example.com", false},
		{4, "another@test.org", false},
		{5, "test@example.com", false},
	}

	cleaned := CleanRecords(records)
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
	"strings"
)

type DataCleaner struct {
	InputPath  string
	OutputPath string
}

func NewDataCleaner(input, output string) *DataCleaner {
	return &DataCleaner{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dc *DataCleaner) Clean() error {
	inputFile, err := os.Open(dc.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dc.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := dc.cleanRow(headers)
	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := dc.cleanRow(record)
		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func (dc *DataCleaner) cleanRow(row []string) []string {
	cleaned := make([]string, len(row))
	for i, value := range row {
		cleaned[i] = strings.TrimSpace(value)
		if cleaned[i] == "" {
			cleaned[i] = "N/A"
		}
	}
	return cleaned
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	cleaner := NewDataCleaner(os.Args[1], os.Args[2])
	if err := cleaner.Clean(); err != nil {
		fmt.Printf("Error cleaning data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data cleaning completed successfully")
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
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

	trimmedHeaders := make([]string, len(headers))
	for i, h := range headers {
		trimmedHeaders[i] = strings.TrimSpace(h)
	}
	if err := writer.Write(trimmedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = strings.TrimSpace(field)
			if cleanedRecord[i] == "" {
				cleanedRecord[i] = "N/A"
			}
		}
		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
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

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data cleaned successfully. Output saved to %s\n", outputFile)
}