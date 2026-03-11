
package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

func sanitizeInput(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func processUserData(data UserData) (UserData, error) {
	data.Username = sanitizeInput(data.Username)
	data.Comments = sanitizeInput(data.Comments)

	if !validateEmail(data.Email) {
		return data, &ValidationError{Field: "email", Message: "invalid email format"}
	}

	if len(data.Username) < 3 || len(data.Username) > 50 {
		return data, &ValidationError{Field: "username", Message: "username must be between 3 and 50 characters"}
	}

	return data, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for i := 0; ; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		record := Record{
			ID:    id,
			Name:  row[1],
			Value: row[2],
		}
		records = append(records, record)
	}

	return records, nil
}

func serializeToJSON(records []Record) (string, error) {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	records, err := processCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		return
	}

	jsonOutput, err := serializeToJSON(records)
	if err != nil {
		fmt.Printf("Error serializing to JSON: %v\n", err)
		return
	}

	fmt.Println(jsonOutput)
}