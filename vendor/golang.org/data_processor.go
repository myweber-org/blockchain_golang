
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

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

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func calculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, r := range records {
		sum += r.Value
		if r.Value > max {
			max = r.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := processCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records\n", len(records))

	avg, max := calculateStats(records)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type CSVProcessor struct {
    filePath   string
    delimiter  rune
    hasHeaders bool
}

func NewCSVProcessor(filePath string, delimiter rune, hasHeaders bool) *CSVProcessor {
    return &CSVProcessor{
        filePath:   filePath,
        delimiter:  delimiter,
        hasHeaders: hasHeaders,
    }
}

func (p *CSVProcessor) ValidateAndClean() ([]map[string]string, error) {
    file, err := os.Open(p.filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = p.delimiter
    reader.TrimLeadingSpace = true

    var headers []string
    var records []map[string]string

    for lineNum := 1; ; lineNum++ {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("line %d: read error: %w", lineNum, err)
        }

        if lineNum == 1 && p.hasHeaders {
            headers = p.cleanRow(row)
            continue
        }

        cleanedRow := p.cleanRow(row)
        if len(cleanedRow) == 0 {
            continue
        }

        var record map[string]string
        if p.hasHeaders {
            if len(headers) != len(cleanedRow) {
                return nil, fmt.Errorf("line %d: column count mismatch", lineNum)
            }
            record = make(map[string]string)
            for i, header := range headers {
                record[header] = cleanedRow[i]
            }
        } else {
            record = make(map[string]string)
            for i, value := range cleanedRow {
                record[fmt.Sprintf("col%d", i+1)] = value
            }
        }
        records = append(records, record)
    }

    return records, nil
}

func (p *CSVProcessor) cleanRow(row []string) []string {
    cleaned := make([]string, 0, len(row))
    for _, cell := range row {
        cleanedCell := strings.TrimSpace(cell)
        if cleanedCell == "" || cleanedCell == "NULL" || cleanedCell == "null" {
            cleanedCell = "N/A"
        }
        cleaned = append(cleaned, cleanedCell)
    }
    return cleaned
}

func main() {
    processor := NewCSVProcessor("input.csv", ',', true)
    records, err := processor.ValidateAndClean()
    if err != nil {
        fmt.Printf("Error processing CSV: %v\n", err)
        return
    }

    fmt.Printf("Successfully processed %d records\n", len(records))
    for i, record := range records {
        if i < 3 {
            fmt.Printf("Record %d: %v\n", i+1, record)
        }
    }
}