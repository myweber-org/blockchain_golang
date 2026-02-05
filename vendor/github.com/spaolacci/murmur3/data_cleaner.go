
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct{}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
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

func (dc *DataCleaner) ProcessRecords(records []string) []string {
	cleaned := dc.RemoveDuplicates(records)
	return cleaned
}

func main() {
	cleaner := &DataCleaner{}
	data := []string{"apple", " banana", "apple", "cherry ", "", "banana", "  "}
	result := cleaner.ProcessRecords(data)
	fmt.Println("Cleaned data:", result)
}
package main

import (
	"fmt"
	"sort"
)

type Record struct {
	ID   int
	Name string
}

type DataSet []Record

func (d DataSet) RemoveDuplicates() DataSet {
	seen := make(map[int]bool)
	result := DataSet{}
	for _, record := range d {
		if !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

func (d DataSet) SortByID() {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})
}

func CleanData(data DataSet) DataSet {
	uniqueData := data.RemoveDuplicates()
	uniqueData.SortByID()
	return uniqueData
}

func main() {
	sampleData := DataSet{
		{3, "Charlie"},
		{1, "Alice"},
		{2, "Bob"},
		{1, "Alice"},
		{4, "David"},
		{2, "Bob"},
	}

	fmt.Println("Original data:", sampleData)
	cleanedData := CleanData(sampleData)
	fmt.Println("Cleaned data:", cleanedData)
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

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := make([]string, len(headers))
	for i, h := range headers {
		cleanedHeaders[i] = strings.TrimSpace(strings.ToLower(h))
	}
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

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			if cleanedField == "" {
				cleanedField = "N/A"
			}
			cleanedRecord[i] = cleanedField
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

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
}