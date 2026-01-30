
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
	ID        int
	Name      string
	Age       int
	Active    bool
	Score     float64
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 5 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 5, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		age, err := strconv.Atoi(strings.TrimSpace(row[2]))
		if err != nil || age < 0 || age > 150 {
			return nil, fmt.Errorf("invalid age at line %d", lineNumber)
		}

		active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
		}

		score, err := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
		if err != nil || score < 0 || score > 100 {
			return nil, fmt.Errorf("invalid score at line %d", lineNumber)
		}

		records = append(records, Record{
			ID:     id,
			Name:   name,
			Age:    age,
			Active: active,
			Score:  score,
		})
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func calculateStatistics(records []Record) (map[string]interface{}, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("empty record set")
	}

	totalAge := 0
	totalScore := 0.0
	activeCount := 0
	highestScore := records[0].Score
	lowestScore := records[0].Score

	for _, record := range records {
		totalAge += record.Age
		totalScore += record.Score

		if record.Active {
			activeCount++
		}

		if record.Score > highestScore {
			highestScore = record.Score
		}
		if record.Score < lowestScore {
			lowestScore = record.Score
		}
	}

	return map[string]interface{}{
		"total_records":   len(records),
		"average_age":     float64(totalAge) / float64(len(records)),
		"average_score":   totalScore / float64(len(records)),
		"active_count":    activeCount,
		"inactive_count":  len(records) - activeCount,
		"highest_score":   highestScore,
		"lowest_score":    lowestScore,
		"active_percentage": float64(activeCount) / float64(len(records)) * 100,
	}, nil
}

func filterRecords(records []Record, filterFunc func(Record) bool) []Record {
	var filtered []Record
	for _, record := range records {
		if filterFunc(record) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	stats, err := calculateStatistics(records)
	if err != nil {
		fmt.Printf("Error calculating statistics: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing completed successfully")
	fmt.Printf("Total records processed: %d\n", stats["total_records"])
	fmt.Printf("Average age: %.2f\n", stats["average_age"])
	fmt.Printf("Average score: %.2f\n", stats["average_score"])
	fmt.Printf("Active users: %d (%.1f%%)\n", stats["active_count"], stats["active_percentage"])

	highScorers := filterRecords(records, func(r Record) bool {
		return r.Score >= 80.0 && r.Active
	})
	fmt.Printf("High performing active users: %d\n", len(highScorers))
}