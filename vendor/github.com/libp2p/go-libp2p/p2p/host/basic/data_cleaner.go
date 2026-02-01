package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func removeDuplicates(inputFile, outputFile string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	reader := csv.NewReader(in)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	var uniqueRecords [][]string

	for _, record := range records {
		if len(record) == 0 {
			continue
		}
		key := record[0]
		if !seen[key] {
			seen[key] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	return writer.WriteAll(uniqueRecords)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := removeDuplicates(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Duplicate removal completed successfully")
}