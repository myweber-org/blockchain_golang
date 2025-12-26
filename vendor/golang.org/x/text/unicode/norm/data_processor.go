package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVData(reader io.Reader) ([]DataRecord, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 4 {
			return nil, errors.New("invalid row format at line " + strconv.Itoa(i+1))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, errors.New("invalid ID at line " + strconv.Itoa(i+1))
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, errors.New("empty name at line " + strconv.Itoa(i+1))
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, errors.New("invalid value at line " + strconv.Itoa(i+1))
		}

		valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		data = append(data, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		})
	}

	return data, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Valid && record.Value > 0 {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateTotalValue(records []DataRecord) float64 {
	var total float64
	for _, record := range records {
		total += record.Value
	}
	return total
}