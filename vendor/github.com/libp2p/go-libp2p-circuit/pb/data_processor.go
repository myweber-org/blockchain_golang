
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

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	headerSkipped := false

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		if len(row) < 3 {
			continue
		}

		age, err := strconv.Atoi(row[1])
		if err != nil {
			continue
		}

		score, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		record := Record{
			Name:  row[0],
			Age:   age,
			Score: score,
		}
		records = append(records, record)
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func calculateAverageScore(records []Record) float64 {
	if len(records) == 0 {
		return 0.0
	}

	total := 0.0
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	jsonOutput, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processed Records:")
	fmt.Println(jsonOutput)

	averageScore := calculateAverageScore(records)
	fmt.Printf("Average Score: %.2f\n", averageScore)
	fmt.Printf("Total Records Processed: %d\n", len(records))
}