package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

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
            return fmt.Errorf("error reading CSV: %w", err)
        }

        cleaned := make([]string, len(record))
        for i, field := range record {
            cleaned[i] = strings.TrimSpace(field)
        }

        if err := writer.Write(cleaned); err != nil {
            return fmt.Errorf("error writing CSV: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run data_cleaner.go <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := cleanCSV(inputFile, outputFile); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully cleaned data from %s to %s\n", inputFile, outputFile)
}
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID      string
	Content string
	Hash    string
}

func generateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		if !seen[record.Hash] {
			seen[record.Hash] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateRecord(record DataRecord) bool {
	if strings.TrimSpace(record.ID) == "" {
		return false
	}
	if strings.TrimSpace(record.Content) == "" {
		return false
	}
	if record.Hash != generateHash(record.Content) {
		return false
	}
	return true
}

func cleanDataPipeline(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if validateRecord(record) {
			validRecords = append(validRecords, record)
		}
	}
	return deduplicateRecords(validRecords)
}

func main() {
	sampleData := []DataRecord{
		{ID: "001", Content: "Sample data record", Hash: generateHash("Sample data record")},
		{ID: "002", Content: "Duplicate record", Hash: generateHash("Duplicate record")},
		{ID: "003", Content: "Sample data record", Hash: generateHash("Sample data record")},
		{ID: "", Content: "Invalid record", Hash: generateHash("Invalid record")},
	}

	cleaned := cleanDataPipeline(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
}