package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	jsonData := []byte(`{"email":"test@example.com","username":"  john_doe  ","age":25}`)
	processedData, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}package main

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
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(line) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(line[0]),
            Name:  strings.TrimSpace(line[1]),
            Email: strings.TrimSpace(line[2]),
            Valid: validateEmail(strings.TrimSpace(line[2])),
        }

        if record.ID != "" && record.Name != "" {
            records = append(records, record)
        }
    }

    return records, nil
}

func validateEmail(email string) bool {
    if email == "" {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func GenerateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid email addresses: %d\n", validCount)
    fmt.Printf("Invalid email addresses: %d\n", len(records)-validCount)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file_path>")
        return
    }

    records, err := ProcessCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    GenerateReport(records)
}
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Valid     bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("record timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := record
	transformed.ID = strings.ToUpper(record.ID)
	transformed.Value = record.Value * 1.1
	transformed.Valid = true

	return transformed, nil
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	var errors []string

	for _, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Record %s: %v", record.ID, err))
			continue
		}
		processed = append(processed, transformed)
	}

	if len(errors) > 0 {
		return processed, fmt.Errorf("processing completed with errors: %s", strings.Join(errors, "; "))
	}

	return processed, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided for statistics calculation")
	}

	var sum float64
	var count int

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			count++
		}
	}

	if count == 0 {
		return 0, 0, errors.New("no valid records found")
	}

	average := sum / float64(count)

	var varianceSum float64
	for _, record := range records {
		if record.Valid {
			diff := record.Value - average
			varianceSum += diff * diff
		}
	}
	variance := varianceSum / float64(count)

	return average, variance, nil
}