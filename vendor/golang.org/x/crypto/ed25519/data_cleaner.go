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
	input := os.Args[1]
	output := os.Args[2]

	if err := cleanCSV(input, output); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully cleaned %s -> %s\n", input, output)
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    result := []DataRecord{}
    for _, record := range records {
        normalizedEmail := strings.ToLower(strings.TrimSpace(record.Email))
        if !seen[normalizedEmail] {
            seen[normalizedEmail] = true
            result = append(result, record)
        }
    }
    return result
}

func ValidateEmails(records []DataRecord) []DataRecord {
    for i := range records {
        email := records[i].Email
        records[i].Valid = strings.Contains(email, "@") && strings.Contains(email, ".")
    }
    return records
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "USER@example.com", false},
        {3, "test.user@domain.org", false},
        {4, "invalid-email", false},
        {5, "user@example.com", false},
    }

    fmt.Println("Original records:", len(sampleData))
    deduped := RemoveDuplicates(sampleData)
    fmt.Println("After deduplication:", len(deduped))
    validated := ValidateEmails(deduped)
    
    for _, record := range validated {
        fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
    }
}