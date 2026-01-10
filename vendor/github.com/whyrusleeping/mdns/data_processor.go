
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		valid, err := strconv.ParseBool(row[3])
		if err != nil {
			return nil, fmt.Errorf("invalid boolean at line %d: %w", lineNumber, err)
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.Valid && record.Value > 0 {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateTotal(records []DataRecord) float64 {
	total := 0.0
	for _, record := range records {
		total += record.Value
	}
	return total
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func NormalizeUserData(data UserData) (UserData, error) {
	if data.Username == "" {
		return data, fmt.Errorf("username cannot be empty")
	}
	if data.Age < 0 || data.Age > 150 {
		return data, fmt.Errorf("invalid age value")
	}

	normalized := UserData{
		Username: strings.TrimSpace(data.Username),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}

	if normalized.Email != "" && !strings.Contains(normalized.Email, "@") {
		return data, fmt.Errorf("invalid email format")
	}

	return normalized, nil
}

func ProcessUserInput(username, email string, age int) {
	data := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	normalized, err := NormalizeUserData(data)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	fmt.Printf("Processed data: %+v\n", normalized)
}

func main() {
	ProcessUserInput("  JohnDoe  ", "JOHN@EXAMPLE.COM", 30)
	ProcessUserInput("", "invalid-email", 200)
}package data

import (
	"errors"
	"regexp"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Tags  []string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeTags(tags []string) []string {
	unique := make(map[string]bool)
	var result []string

	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !unique[normalized] {
			unique[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func ProcessRecord(record Record) (Record, error) {
	if err := ValidateEmail(record.Email); err != nil {
		return Record{}, err
	}

	processed := Record{
		ID:    strings.TrimSpace(record.ID),
		Email: strings.ToLower(strings.TrimSpace(record.Email)),
		Tags:  NormalizeTags(record.Tags),
	}

	if processed.ID == "" {
		return Record{}, errors.New("record ID cannot be empty")
	}

	return processed, nil
}