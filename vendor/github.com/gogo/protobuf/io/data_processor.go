
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Email   string
	Active  bool
	Score   float64
}

func ParseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 5 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 5, got %d", lineNum, len(row))
		}

		record, err := parseRow(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("no valid records found in file")
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNum)
	}

	record.Email = strings.TrimSpace(row[2])
	if !strings.Contains(record.Email, "@") {
		return record, fmt.Errorf("invalid email format at line %d", lineNum)
	}

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}
	record.Active = active

	score, err := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid score at line %d: %w", lineNum, err)
	}
	record.Score = score

	return record, nil
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	seenEmails := make(map[string]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if seenEmails[record.Email] {
			return fmt.Errorf("duplicate email found: %s", record.Email)
		}
		seenEmails[record.Email] = true

		if record.Score < 0 || record.Score > 100 {
			return fmt.Errorf("score out of range for record %d: %f", record.ID, record.Score)
		}
	}

	return nil
}

func CalculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	var maxScore float64

	for _, record := range records {
		sum += record.Score
		if record.Active {
			activeCount++
		}
		if record.Score > maxScore {
			maxScore = record.Score
		}
	}

	averageScore := sum / float64(len(records))
	return averageScore, maxScore, activeCount
}