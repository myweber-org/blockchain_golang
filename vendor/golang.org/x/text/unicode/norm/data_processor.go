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
}package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	userData := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	userData = TransformUsername(userData)

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}