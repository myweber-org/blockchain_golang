package main

import (
    "encoding/csv"
    "errors"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
    Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) < 4 {
            continue
        }

        record, err := validateAndCreateRecord(row)
        if err != nil {
            continue
        }

        records = append(records, record)
    }

    return records, nil
}

func validateAndCreateRecord(row []string) (DataRecord, error) {
    var record DataRecord

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, errors.New("invalid id")
    }

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, errors.New("empty name")
    }

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, errors.New("invalid value")
    }

    valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, errors.New("invalid valid flag")
    }

    record.ID = id
    record.Name = name
    record.Value = value
    record.Valid = valid

    return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Valid {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func CalculateAverage(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0.0
    }

    total := 0.0
    for _, record := range records {
        total += record.Value
    }

    return total / float64(len(records))
}
package main

import (
	"regexp"
	"strings"
)

func ProcessInput(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	trimmed := strings.TrimSpace(input)

	re := regexp.MustCompile(`[^a-zA-Z0-9\s\-_]`)
	cleaned := re.ReplaceAllString(trimmed, "")

	if len(cleaned) > 100 {
		cleaned = cleaned[:100]
	}

	return cleaned, nil
}
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) != 3 {
			return nil, errors.New("invalid row length at line " + strconv.Itoa(i+1))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid ID format at line " + strconv.Itoa(i+1))
		}

		name := row[1]
		if name == "" {
			return nil, errors.New("empty name at line " + strconv.Itoa(i+1))
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value format at line " + strconv.Itoa(i+1))
		}

		data = append(data, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return data, nil
}

func ValidateData(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid ID: " + strconv.Itoa(record.ID))
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate ID: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true
	}
	return nil
}

func WriteProcessedData(outputPath string, records []DataRecord) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			strconv.FormatFloat(record.Value, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func ProcessCSVData(inputPath, outputPath string) error {
	records, err := ParseCSVFile(inputPath)
	if err != nil {
		return err
	}

	if err := ValidateData(records); err != nil {
		return err
	}

	return WriteProcessedData(outputPath, records)
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Biography string `json:"biography,omitempty"`
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", profile.ID)
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}

	return nil
}

func TransformUserProfile(profile *UserProfile) {
	profile.Username = strings.TrimSpace(profile.Username)
	profile.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	profile.Biography = strings.TrimSpace(profile.Biography)

	if len(profile.Biography) > 500 {
		profile.Biography = profile.Biography[:500] + "..."
	}
}

func ProcessUserProfile(data []byte) ([]byte, error) {
	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user profile: %w", err)
	}

	if err := ValidateUserProfile(profile); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	TransformUserProfile(&profile)

	processedData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal processed profile: %w", err)
	}

	return processedData, nil
}

func main() {
	jsonData := []byte(`{
		"id": 123,
		"username": "  JohnDoe  ",
		"email": "JOHN@EXAMPLE.COM",
		"age": 30,
		"active": true,
		"biography": "This is a very long biography that exceeds the maximum allowed length. It contains many details about the user's life, interests, and professional background. We need to ensure it gets properly truncated during processing to maintain data consistency across our systems."
	}`)

	result, err := ProcessUserProfile(jsonData)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	fmt.Println("Processed user profile:")
	fmt.Println(string(result))
}