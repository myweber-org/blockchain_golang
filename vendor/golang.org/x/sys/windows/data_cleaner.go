
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
	
	data := []string{"Apple", "apple", " BANANA ", "banana", "Cherry"}
	
	fmt.Println("Original data:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated:", deduped)
	
	cleaner.Reset()
	
	testValue := "  TEST  "
	fmt.Printf("Normalized '%s': '%s'\n", testValue, cleaner.Normalize(testValue))
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func deduplicateRecords(records [][]string) [][]string {
	seen := make(map[string]bool)
	var unique [][]string
	for _, record := range records {
		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func normalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func normalizeRecords(records [][]string) [][]string {
	for i := range records {
		for j := range records[i] {
			records[i][j] = normalizeString(records[i][j])
		}
	}
	return records
}

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func writeCSV(filename string, records [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	return writer.WriteAll(records)
}

func sortRecords(records [][]string, columnIndex int) {
	if len(records) == 0 || columnIndex < 0 {
		return
	}
	sort.Slice(records, func(i, j int) bool {
		if columnIndex >= len(records[i]) || columnIndex >= len(records[j]) {
			return false
		}
		return records[i][columnIndex] < records[j][columnIndex]
	})
}

func processData(inputFile, outputFile string) error {
	records, err := readCSV(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no data found in input file")
	}

	records = deduplicateRecords(records)
	records = normalizeRecords(records)
	sortRecords(records, 0)

	if err := writeCSV(outputFile, records); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Processed %d records, saved to %s\n", len(records), outputFile)
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}