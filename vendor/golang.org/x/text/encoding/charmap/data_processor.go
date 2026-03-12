package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func SanitizeInput(input string) string {
	var result strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Username = NormalizeUsername(data.Username)
	data.Username = SanitizeInput(data.Username)

	if !ValidateEmail(data.Email) {
		return data, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return data, fmt.Errorf("age out of valid range")
	}

	return data, nil
}

func main() {
	sampleData := UserData{
		Username: "  John_Doe123  ",
		Email:    "john@example.com",
		Age:      30,
	}

	processed, err := ProcessUserData(sampleData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed user: %+v\n", processed)
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

func FormatJSON(input string) (string, error) {
    var data interface{}
    err := json.Unmarshal([]byte(input), &data)
    if err != nil {
        return "", fmt.Errorf("invalid JSON: %w", err)
    }

    formatted, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to format JSON: %w", err)
    }

    return string(formatted), nil
}

func ValidateJSON(input string) bool {
    var js json.RawMessage
    return json.Unmarshal([]byte(input), &js) == nil
}

func main() {
    sample := `{"name":"test","value":123,"active":true}`
    fmt.Println("Is valid?", ValidateJSON(sample))

    formatted, err := FormatJSON(sample)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Formatted JSON:")
    fmt.Println(formatted)
}
package main

import (
	"encoding/csv"
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

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
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

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNum)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}
	record.Active = active

	return record, nil
}

func validateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value not allowed for record %d", record.ID)
		}
	}

	return nil
}

func processData(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	totalValue := 0.0
	activeCount := 0

	for _, record := range records {
		totalValue += record.Value
		if record.Active {
			activeCount++
		}
	}

	fmt.Printf("Processed %d records successfully\n", len(records))
	fmt.Printf("Total value: %.2f\n", totalValue)
	fmt.Printf("Active records: %d\n", activeCount)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if err := processData(filename); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}