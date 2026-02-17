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

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleaned); err != nil {
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

	fmt.Printf("Successfully cleaned CSV file: %s -> %s\n", inputFile, outputFile)
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    var unique []DataRecord

    for _, record := range records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
        if !seen[key] {
            seen[key] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidatePhone(phone string) bool {
    return len(phone) >= 10 && strings.Count(phone, "") > 1
}

func CleanData(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    for _, record := range records {
        if ValidateEmail(record.Email) && ValidatePhone(record.Phone) {
            validRecords = append(validRecords, record)
        }
    }
    return DeduplicateRecords(validRecords)
}

func main() {
    sampleData := []DataRecord{
        {1, "test@example.com", "1234567890"},
        {2, "invalid-email", "5555555555"},
        {3, "test@example.com", "1234567890"},
        {4, "another@test.org", "0987654321"},
    }

    cleaned := CleanData(sampleData)
    fmt.Printf("Original: %d records\n", len(sampleData))
    fmt.Printf("Cleaned: %d records\n", len(cleaned))
}package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Cleaned:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Cleaned:", uniqueStrings)
}