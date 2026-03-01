
package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", data.ID)
	}
	if strings.TrimSpace(data.Name) == "" {
		return fmt.Errorf("user name cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format: %s", data.Email)
	}
	return nil
}

func TransformUserName(data UserData) UserData {
	data.Name = strings.ToUpper(strings.TrimSpace(data.Name))
	return data
}

func ProcessUserData(data UserData) (UserData, error) {
	if err := ValidateUserData(data); err != nil {
		return UserData{}, err
	}
	return TransformUserName(data), nil
}

func main() {
	user := UserData{
		ID:    1001,
		Name:  "  john doe  ",
		Email: "john@example.com",
	}

	processed, err := ProcessUserData(user)
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		return
	}

	fmt.Printf("Processed user: %+v\n", processed)
}
package main

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("username can only contain letters, digits, underscores, and hyphens")
		}
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TransformUserData(username, email string) (string, string, error) {
	if err := ValidateUsername(username); err != nil {
		return "", "", err
	}

	normalizedEmail := NormalizeEmail(email)
	return username, normalizedEmail, nil
}
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
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format at line %d: %w", lineNumber, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNumber++
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file_path>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	records, err := ProcessCSVFile(filePath)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	average, max := CalculateStatistics(records)
	fmt.Printf("Processed %d records successfully\n", len(records))
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Maximum value: %.2f\n", max)
}
package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Value   float64 `json:"value"`
    Active  bool    `json:"active"`
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
            return nil, fmt.Errorf("csv read error on line %d: %w", lineNumber, err)
        }

        if len(row) != 4 {
            return nil, fmt.Errorf("invalid column count on line %d: expected 4, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID on line %d: %w", lineNumber, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value on line %d: %w", lineNumber, err)
        }

        active, err := strconv.ParseBool(row[3])
        if err != nil {
            return nil, fmt.Errorf("invalid active flag on line %d: %w", lineNumber, err)
        }

        record := Record{
            ID:     id,
            Name:   row[1],
            Value:  value,
            Active: active,
        }
        records = append(records, record)
    }

    return records, nil
}

func convertToJSON(records []Record) (string, error) {
    jsonData, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return "", fmt.Errorf("json marshaling failed: %w", err)
    }
    return string(jsonData), nil
}

func filterActiveRecords(records []Record) []Record {
    var activeRecords []Record
    for _, record := range records {
        if record.Active {
            activeRecords = append(activeRecords, record)
        }
    }
    return activeRecords
}

func calculateTotalValue(records []Record) float64 {
    var total float64
    for _, record := range records {
        total += record.Value
    }
    return total
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

    fmt.Printf("Processed %d records\n", len(records))

    activeRecords := filterActiveRecords(records)
    fmt.Printf("Active records: %d\n", len(activeRecords))

    totalValue := calculateTotalValue(records)
    fmt.Printf("Total value: %.2f\n", totalValue)

    jsonOutput, err := convertToJSON(records)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        os.Exit(1)
    }

    outputFile := "output.json"
    err = os.WriteFile(outputFile, []byte(jsonOutput), 0644)
    if err != nil {
        fmt.Printf("Error writing JSON file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("JSON output written to %s\n", outputFile)
}