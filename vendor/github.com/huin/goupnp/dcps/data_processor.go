
package main

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func TransformPhoneNumber(phone string) (string, error) {
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	if len(cleaned) < 10 {
		return "", errors.New("phone number too short")
	}
	if len(cleaned) > 15 {
		return "", errors.New("phone number too long")
	}
	return "+" + cleaned, nil
}

func ProcessUserData(email, username, phone string) (map[string]string, error) {
	result := make(map[string]string)

	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	result["email"] = email

	normalizedUser := NormalizeUsername(username)
	result["username"] = normalizedUser

	transformedPhone, err := TransformPhoneNumber(phone)
	if err != nil {
		return nil, err
	}
	result["phone"] = transformedPhone

	return result, nil
}
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

type DataRecord struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNum, len(row))
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

func parseRow(row []string, lineNum int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}
	record.Active = active

	return record, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d (must be positive)", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value not allowed for record ID %d", record.ID)
		}
	}

	return nil
}

func CalculateStats(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	var maxValue float64

	for _, record := range records {
		sum += record.Value
		if record.Value > maxValue {
			maxValue = record.Value
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, maxValue, activeCount
}