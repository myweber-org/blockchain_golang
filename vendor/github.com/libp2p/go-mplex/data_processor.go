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

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data *UserData) {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	user := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	TransformUsername(&user)

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	return user, nil
}

func main() {
	user, err := ProcessUserInput("  JohnDoe  ", "john@example.com", 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}