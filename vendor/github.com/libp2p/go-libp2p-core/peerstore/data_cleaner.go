
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataCleaner struct {
    inputPath  string
    outputPath string
    seenRows   map[string]bool
}

func NewDataCleaner(input, output string) *DataCleaner {
    return &DataCleaner{
        inputPath:  input,
        outputPath: output,
        seenRows:   make(map[string]bool),
    }
}

func (dc *DataCleaner) Clean() error {
    inputFile, err := os.Open(dc.inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(dc.outputPath)
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

        rowKey := strings.Join(record, "|")
        if !dc.seenRows[rowKey] {
            dc.seenRows[rowKey] = true
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("failed to write record: %w", err)
            }
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    cleaner := NewDataCleaner(os.Args[1], os.Args[2])
    if err := cleaner.Clean(); err != nil {
        fmt.Printf("Error cleaning data: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Data cleaned successfully. Output written to %s\n", os.Args[2])
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID   string
	Name string
	Email string
}

func generateHash(record DataRecord) string {
	data := fmt.Sprintf("%s|%s|%s", record.ID, record.Name, record.Email)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	
	for _, record := range records {
		hash := generateHash(record)
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func cleanData(records []DataRecord) []DataRecord {
	var valid []DataRecord
	
	for _, record := range records {
		if validateEmail(record.Email) {
			valid = append(valid, record)
		}
	}
	
	return deduplicateRecords(valid)
}

func main() {
	records := []DataRecord{
		{"1", "John Doe", "john@example.com"},
		{"2", "Jane Smith", "jane@example.com"},
		{"3", "John Doe", "john@example.com"},
		{"4", "Bob Wilson", "invalid-email"},
		{"5", "Jane Smith", "jane@example.com"},
	}
	
	cleaned := cleanData(records)
	
	fmt.Printf("Original records: %d\n", len(records))
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	
	for _, record := range cleaned {
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", 
			record.ID, record.Name, record.Email)
	}
}