
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Score float64
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNum, err)
		}

		if len(row) != 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])
		email := strings.TrimSpace(row[2])
		score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			continue
		}

		if !isValidEmail(email) {
			continue
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		})
	}

	return records, nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func calculateAverageScore(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	total := 0.0
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func filterByScoreThreshold(records []DataRecord, threshold float64) []DataRecord {
	filtered := []DataRecord{}
	for _, record := range records {
		if record.Score >= threshold {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <csv_file>")
		return
	}

	records, err := parseCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d valid records\n", len(records))
	
	avgScore := calculateAverageScore(records)
	fmt.Printf("Average score: %.2f\n", avgScore)

	highScorers := filterByScoreThreshold(records, 80.0)
	fmt.Printf("Records with score >= 80: %d\n", len(highScorers))
}