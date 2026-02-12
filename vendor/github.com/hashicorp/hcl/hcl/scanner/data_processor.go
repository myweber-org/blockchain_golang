package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}
package main

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
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 3 {
			return nil, errors.New("invalid csv format at line " + strconv.Itoa(i+1))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid id at line " + strconv.Itoa(i+1))
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value at line " + strconv.Itoa(i+1))
		}

		data = append(data, DataRecord{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return data, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid id: " + strconv.Itoa(record.ID))
		}
		if record.Name == "" {
			return errors.New("empty name for id: " + strconv.Itoa(record.ID))
		}
		if record.Value < 0 {
			return errors.New("negative value for id: " + strconv.Itoa(record.ID))
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate id: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true
	}
	return nil
}

func WriteProcessedData(outputPath string, records []DataRecord) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Value", "Processed"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		processedValue := record.Value * 1.1
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			strconv.FormatFloat(record.Value, 'f', 2, 64),
			strconv.FormatFloat(processedValue, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}