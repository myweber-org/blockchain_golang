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
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

func processCSVData(filename string) ([]Record, error) {
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

		if len(row) < 3 {
			continue
		}

		value, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			continue
		}

		count, err := strconv.Atoi(row[2])
		if err != nil {
			continue
		}

		records = append(records, Record{
			Name:  row[0],
			Value: value,
			Count: count,
		})
	}

	return records, nil
}

func generateJSONReport(records []Record) (string, error) {
	report := struct {
		TotalRecords int       `json:"total_records"`
		AverageValue float64   `json:"average_value"`
		MaxCount     int       `json:"max_count"`
		Records      []Record  `json:"records"`
	}{
		TotalRecords: len(records),
		Records:      records,
	}

	if len(records) > 0 {
		var totalValue float64
		maxCount := 0

		for _, r := range records {
			totalValue += r.Value
			if r.Count > maxCount {
				maxCount = r.Count
			}
		}

		report.AverageValue = totalValue / float64(len(records))
		report.MaxCount = maxCount
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := processCSVData(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	jsonReport, err := generateJSONReport(records)
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(jsonReport)
}