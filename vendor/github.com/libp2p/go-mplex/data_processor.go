package main

import (
	"encoding/csv"
	"errors"
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

func ParseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]Record, 0)

	// Skip header
	if _, err := reader.Read(); err != nil {
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

		if len(row) != 4 {
			return nil, errors.New("invalid row length")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		name := strings.TrimSpace(row[1])
		email := strings.TrimSpace(row[2])
		score, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, err
		}

		record := Record{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		}

		if !validateRecord(record) {
			return nil, errors.New("invalid record data")
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(r Record) bool {
	if r.ID <= 0 {
		return false
	}
	if r.Name == "" {
		return false
	}
	if !strings.Contains(r.Email, "@") {
		return false
	}
	if r.Score < 0 || r.Score > 100 {
		return false
	}
	return true
}

func CalculateAverageScore(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}

	var total float64
	for _, r := range records {
		total += r.Score
	}

	return total / float64(len(records))
}