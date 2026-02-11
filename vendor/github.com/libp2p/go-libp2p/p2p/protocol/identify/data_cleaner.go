
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Score int
}

func readCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(row) < 3 {
			continue
		}

		record := Record{
			ID:    strings.TrimSpace(row[0]),
			Email: strings.TrimSpace(row[1]),
			Score: 0,
		}
		fmt.Sscanf(row[2], "%d", &record.Score)
		records = append(records, record)
	}
	return records, nil
}

func deduplicate(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		key := record.ID + "|" + record.Email
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateRecords(records []Record) []Record {
	var valid []Record
	for _, record := range records {
		if record.ID != "" && strings.Contains(record.Email, "@") && record.Score >= 0 {
			valid = append(valid, record)
		}
	}
	return valid
}

func sortByScore(records []Record) []Record {
	sort.Slice(records, func(i, j int) bool {
		return records[i].Score > records[j].Score
	})
	return records
}

func writeCSV(filename string, records []Record) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Email", "Score"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		row := []string{record.ID, record.Email, fmt.Sprintf("%d", record.Score)}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	records, err := readCSV("input.csv")
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	records = deduplicate(records)
	records = validateRecords(records)
	records = sortByScore(records)

	if err := writeCSV("cleaned_data.csv", records); err != nil {
		fmt.Printf("Error writing CSV: %v\n", err)
		return
	}

	fmt.Printf("Processed %d valid records\n", len(records))
}