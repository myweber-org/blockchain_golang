
package main

import (
    "fmt"
    "strings"
    "unicode"
)

type UserData struct {
    Username string
    Email    string
}

func NormalizeUsername(username string) string {
    trimmed := strings.TrimSpace(username)
    var result strings.Builder
    for _, r := range trimmed {
        if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
            result.WriteRune(unicode.ToLower(r))
        }
    }
    return result.String()
}

func ValidateEmail(email string) bool {
    trimmed := strings.TrimSpace(email)
    if len(trimmed) < 3 || len(trimmed) > 254 {
        return false
    }
    atIndex := strings.LastIndex(trimmed, "@")
    if atIndex < 1 || atIndex == len(trimmed)-1 {
        return false
    }
    dotIndex := strings.LastIndex(trimmed[atIndex:], ".")
    if dotIndex < 1 || dotIndex == len(trimmed[atIndex:])-1 {
        return false
    }
    return true
}

func ProcessUserInput(username, email string) (*UserData, error) {
    normalizedUsername := NormalizeUsername(username)
    if normalizedUsername == "" {
        return nil, fmt.Errorf("invalid username: contains no valid characters")
    }

    if !ValidateEmail(email) {
        return nil, fmt.Errorf("invalid email format")
    }

    return &UserData{
        Username: normalizedUsername,
        Email:    strings.ToLower(strings.TrimSpace(email)),
    }, nil
}

func main() {
    testData := []struct {
        username string
        email    string
    }{
        {"  John_Doe-123  ", "john@example.com"},
        {"Alice.Bob", "alice@test.org"},
        {"   ", "invalid-email"},
        {"Test-User_1", "bad@email"},
    }

    for _, td := range testData {
        user, err := ProcessUserInput(td.username, td.email)
        if err != nil {
            fmt.Printf("Error processing %s, %s: %v\n", td.username, td.email, err)
        } else {
            fmt.Printf("Processed: %+v\n", user)
        }
    }
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
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func processCSVFile(filename string) ([]Record, error) {
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRecord(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRecord(row []string, lineNumber int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
	}
	record.Active = active

	return record, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Name == "" {
			return fmt.Errorf("empty name for record ID: %d", record.ID)
		}

		if record.Value < 0 {
			return fmt.Errorf("negative value for record ID: %d", record.ID)
		}
	}

	return nil
}

func calculateStatistics(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	minValue := records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value < minValue {
			minValue = record.Value
		}
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, minValue, activeCount
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	if err := validateRecords(records); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	average, minValue, activeCount := calculateStatistics(records)

	fmt.Printf("Processing completed successfully\n")
	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Minimum value: %.2f\n", minValue)
	fmt.Printf("Active records: %d\n", activeCount)
}