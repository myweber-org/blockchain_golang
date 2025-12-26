package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID    int
	Name  string
	Email string
	Score float64
}

func parseCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNum++
		if lineNum == 1 {
			continue
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

		if !validateEmail(email) {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		})
	}

	return records, nil
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func calculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Score

	for _, r := range records {
		sum += r.Score
		if r.Score > max {
			max = r.Score
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <csv_file>")
		return
	}

	records, err := parseCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	fmt.Printf("Processed %d valid records\n", len(records))
	
	avg, max := calculateStats(records)
	fmt.Printf("Average score: %.2f\n", avg)
	fmt.Printf("Maximum score: %.2f\n", max)

	for i, r := range records {
		if i < 3 {
			fmt.Printf("Sample record: ID=%d, Name=%s, Email=%s\n", r.ID, r.Name, r.Email)
		}
	}
}