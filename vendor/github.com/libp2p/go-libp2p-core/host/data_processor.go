package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to unmarshal JSON data into a User struct.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ParseUserData(rawData []byte) (*User, error) {
	valid, err := ValidateJSON(rawData)
	if !valid {
		return nil, err
	}

	var user User
	err = json.Unmarshal(rawData, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func main() {
	jsonData := []byte(`{"id": 1, "name": "Alice", "email": "alice@example.com"}`)

	user, err := ParseUserData(jsonData)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Parsed User: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
}
package main

import (
    "encoding/csv"
    "errors"
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

func ParseCSVFile(filename string) ([]DataRecord, error) {
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

        if lineNumber == 1 {
            continue
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

func parseRow(row []string, lineNumber int) (DataRecord, error) {
    var record DataRecord

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }
    record.Value = value

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
    }
    record.Active = active

    return record, nil
}

func ValidateRecords(records []DataRecord) error {
    if len(records) == 0 {
        return errors.New("no records to validate")
    }

    idSet := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid ID %d: must be positive", record.ID)
        }

        if idSet[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        idSet[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value for record ID %d", record.ID)
        }
    }

    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var activeCount int
    var minValue float64 = records[0].Value

    for i, record := range records {
        sum += record.Value
        if record.Active {
            activeCount++
        }
        if i == 0 || record.Value < minValue {
            minValue = record.Value
        }
    }

    average := sum / float64(len(records))
    return average, minValue, activeCount
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
    var activeRecords []DataRecord
    for _, record := range records {
        if record.Active {
            activeRecords = append(activeRecords, record)
        }
    }
    return activeRecords
}