
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
}

func processCSVFile(inputPath string) ([]Record, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", lineNumber)
		}

		age, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, fmt.Errorf("invalid age at line %d: %w", lineNumber, err)
		}

		score, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid score at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			Name:  line[0],
			Age:   age,
			Score: score,
		})
		lineNumber++
	}

	return records, nil
}

func convertToJSON(records []Record, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("json encoding failed: %w", err)
	}

	return nil
}

func calculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var totalScore float64
	minScore := records[0].Score
	maxScore := records[0].Score

	for _, record := range records {
		totalScore += record.Score
		if record.Score < minScore {
			minScore = record.Score
		}
		if record.Score > maxScore {
			maxScore = record.Score
		}
	}

	averageScore := totalScore / float64(len(records))
	return averageScore, minScore, int(maxScore)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.json>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	records, err := processCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	if err := convertToJSON(records, outputFile); err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	avgScore, minScore, maxScore := calculateStats(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Average score: %.2f\n", avgScore)
	fmt.Printf("Score range: %.2f - %d\n", minScore, maxScore)
	fmt.Printf("JSON output saved to: %s\n", outputFile)
}