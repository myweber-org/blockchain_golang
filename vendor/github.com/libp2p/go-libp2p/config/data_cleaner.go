package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, rec := range records {
		email := strings.ToLower(strings.TrimSpace(rec.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, rec)
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func CleanData(records []Record) ([]Record, error) {
	var cleaned []Record
	for _, rec := range records {
		if err := ValidateEmail(rec.Email); err != nil {
			continue
		}
		cleaned = append(cleaned, rec)
	}
	cleaned = DeduplicateRecords(cleaned)
	if len(cleaned) == 0 {
		return cleaned, errors.New("no valid records after cleaning")
	}
	return cleaned, nil
}

func main() {
	sampleData := []Record{
		{1, "user@example.com", true},
		{2, "invalid-email", true},
		{3, "user@example.com", true},
		{4, "another@test.org", true},
		{5, "", true},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Cleaned %d records\n", len(cleaned))
	for _, rec := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", rec.ID, rec.Email)
	}
}
package main

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
    Delimiter  rune
}

func NewDataCleaner(input, output string) *DataCleaner {
    return &DataCleaner{
        InputPath:  input,
        OutputPath: output,
        Delimiter:  ',',
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
    reader.Comma = dc.Delimiter

    writer := csv.NewWriter(outputFile)
    writer.Comma = dc.Delimiter
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    seen := make(map[string]bool)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read record: %w", err)
        }

        key := strings.Join(record, "|")
        if seen[key] {
            continue
        }
        seen[key] = true

        for i, value := range record {
            record[i] = strings.TrimSpace(value)
            if record[i] == "" {
                record[i] = "N/A"
            }
        }

        if err := writer.Write(record); err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    cleaner := NewDataCleaner(os.Args[1], os.Args[2])
    if err := cleaner.Clean(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}