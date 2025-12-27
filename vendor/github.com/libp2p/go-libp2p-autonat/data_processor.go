package data_processor

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) != 3 {
			return nil, errors.New("invalid column count at line " + strconv.Itoa(lineNumber))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, errors.New("invalid ID at line " + strconv.Itoa(lineNumber))
		}

		name := line[1]
		if name == "" {
			return nil, errors.New("empty name at line " + strconv.Itoa(lineNumber))
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, errors.New("invalid value at line " + strconv.Itoa(lineNumber))
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []DataRecord) {
	valid := []DataRecord{}
	invalid := []DataRecord{}

	for _, record := range records {
		if record.ID > 0 && record.Value >= 0 {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}

	return valid, invalid
}