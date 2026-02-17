
package main

import (
	"fmt"
	"sort"
)

func FilterAndSort(numbers []int, threshold int) []int {
	var filtered []int
	for _, num := range numbers {
		if num > threshold {
			filtered = append(filtered, num)
		}
	}
	sort.Ints(filtered)
	return filtered
}

func main() {
	data := []int{45, 12, 89, 3, 67, 34, 9, 21}
	result := FilterAndSort(data, 20)
	fmt.Println("Filtered and sorted:", result)
}
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

func ProcessCSVFile(filepath string) ([]DataRecord, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(line[0]),
			Name:  strings.TrimSpace(line[1]),
			Email: strings.TrimSpace(line[2]),
			Valid: validateRecord(strings.TrimSpace(line[0]), strings.TrimSpace(line[2])),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file_path>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	validRecords := FilterValidRecords(records)
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))

	for _, record := range validRecords {
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}