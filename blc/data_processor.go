
package main

import (
	"errors"
	"regexp"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Score int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if !emailRegex.MatchString(record.Email) {
		return errors.New("invalid email format")
	}
	if record.Score < 0 || record.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ProcessRecords(records []DataRecord) ([]DataRecord, []error) {
	var processed []DataRecord
	var errs []error

	for i, record := range records {
		record.Email = NormalizeEmail(record.Email)
		if err := ValidateRecord(record); err != nil {
			errs = append(errs, errors.New("record "+record.ID+": "+err.Error()))
			continue
		}
		processed = append(processed, record)
		_ = i
	}

	return processed, errs
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}
	var sum int
	for _, record := range records {
		sum += record.Score
	}
	return float64(sum) / float64(len(records))
}