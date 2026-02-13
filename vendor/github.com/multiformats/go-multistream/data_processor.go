
package main

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserInput(email, username string, age int) (*UserData, error) {
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	sanitizedUsername := SanitizeUsername(username)
	if sanitizedUsername == "" {
		return nil, errors.New("username cannot be empty")
	}

	if err := ValidateAge(age); err != nil {
		return nil, err
	}

	return &UserData{
		Email:    email,
		Username: sanitizedUsername,
		Age:      age,
	}, nil
}
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

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	line := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", line, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", line)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
		line++
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, r := range records {
		if r.ID <= 0 {
			return fmt.Errorf("invalid ID %d", r.ID)
		}
		if seen[r.ID] {
			return fmt.Errorf("duplicate ID %d", r.ID)
		}
		if r.Value < 0 {
			return fmt.Errorf("negative value for ID %d", r.ID)
		}
		seen[r.ID] = true
	}
	return nil
}

func CalculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}