package main

import (
	"errors"
	"strings"
	"unicode"
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
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range data.Username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username contains invalid characters")
		}
	}

	if !strings.Contains(data.Email, "@") || !strings.Contains(data.Email, ".") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func TransformEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	localPart := strings.ToLower(parts[0])
	domain := strings.ToLower(parts[1])
	return localPart + "@" + domain
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID    string
	Name  string
	Email string
	Valid bool
}

func ProcessCSVFile(filepath string) ([]DataRecord, error) {
	file, err := os.Open(filepath)
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

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(strings.TrimSpace(row[0]), strings.TrimSpace(row[2])),
		}

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	validRecords := []DataRecord{}
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	validCount := 0
	for _, record := range records {
		if record.Valid {
			validCount++
		}
	}
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", len(records)-validCount)
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

type Record struct {
	ID        int
	Name      string
	Age       int
	Active    bool
	Score     float64
	Timestamp string
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

		if len(row) != 6 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 6, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (Record, error) {
	var record Record
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNumber)
	}

	record.Age, err = strconv.Atoi(strings.TrimSpace(row[2]))
	if err != nil || record.Age < 0 || record.Age > 150 {
		return Record{}, fmt.Errorf("invalid age at line %d", lineNumber)
	}

	record.Active, err = strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}

	record.Score, err = strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
	if err != nil || record.Score < 0 || record.Score > 100 {
		return Record{}, fmt.Errorf("invalid score at line %d", lineNumber)
	}

	record.Timestamp = strings.TrimSpace(row[5])
	if record.Timestamp == "" {
		return Record{}, fmt.Errorf("empty timestamp at line %d", lineNumber)
	}

	return record, nil
}

func validateRecords(records []Record) error {
	if len(records) == 0 {
		return fmt.Errorf("no records found")
	}

	idSet := make(map[int]bool)
	for _, record := range records {
		if idSet[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		idSet[record.ID] = true
	}

	return nil
}

func calculateStatistics(records []Record) (map[string]interface{}, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("cannot calculate statistics for empty records")
	}

	stats := make(map[string]interface{})
	var totalScore float64
	var totalAge int
	activeCount := 0

	for _, record := range records {
		totalScore += record.Score
		totalAge += record.Age
		if record.Active {
			activeCount++
		}
	}

	stats["total_records"] = len(records)
	stats["average_score"] = totalScore / float64(len(records))
	stats["average_age"] = float64(totalAge) / float64(len(records))
	stats["active_percentage"] = (float64(activeCount) / float64(len(records))) * 100
	stats["active_count"] = activeCount

	return stats, nil
}

func processDataFile(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	stats, err := calculateStatistics(records)
	if err != nil {
		return fmt.Errorf("statistics calculation failed: %w", err)
	}

	fmt.Println("Data processing completed successfully")
	fmt.Printf("Processed %d records\n", stats["total_records"])
	fmt.Printf("Average Score: %.2f\n", stats["average_score"])
	fmt.Printf("Average Age: %.2f\n", stats["average_age"])
	fmt.Printf("Active Users: %d (%.1f%%)\n", stats["active_count"], stats["active_percentage"])

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if err := processDataFile(filename); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}