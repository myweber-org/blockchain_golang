
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
}

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return fmt.Errorf("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 120 {
		return fmt.Errorf("age must be between 0 and 120")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(transformed.Username)
	transformed.Email = strings.TrimSpace(transformed.Email)

	if transformed.Age < 18 {
		transformed.Active = false
	}

	return transformed
}

func ProcessUserData(input []byte) ([]byte, error) {
	var profile UserProfile
	if err := json.Unmarshal(input, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if err := ValidateUserProfile(profile); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	transformed := TransformProfile(profile)
	output, err := json.MarshalIndent(transformed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return output, nil
}

func main() {
	inputData := []byte(`{
		"id": 123,
		"username": "JohnDoe",
		"email": "JOHN@EXAMPLE.COM",
		"age": 25,
		"active": true
	}`)

	result, err := ProcessUserData(inputData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Processed user profile:")
	fmt.Println(string(result))
}
package main

import "fmt"

func movingAverage(data []float64, windowSize int) []float64 {
    if len(data) == 0 || windowSize <= 0 || windowSize > len(data) {
        return []float64{}
    }

    result := make([]float64, len(data)-windowSize+1)
    var sum float64

    for i := 0; i < windowSize; i++ {
        sum += data[i]
    }
    result[0] = sum / float64(windowSize)

    for i := windowSize; i < len(data); i++ {
        sum = sum - data[i-windowSize] + data[i]
        result[i-windowSize+1] = sum / float64(windowSize)
    }

    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := movingAverage(sampleData, window)
    fmt.Printf("Moving averages with window %d: %v\n", window, averages)
}
package main

import "fmt"

func calculateAverage(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    
    return float64(sum) / float64(len(numbers))
}

func main() {
    data := []int{10, 20, 30, 40, 50}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
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
    ID      string
    Name    string
    Email   string
    Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
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

        if len(row) < 4 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if record.ID == "" || record.Name == "" {
            return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
        }

        if !strings.Contains(record.Email, "@") {
            return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    for _, record := range records {
        if record.Active == "true" && strings.HasSuffix(record.Email, ".com") {
            validRecords = append(validRecords, record)
        }
    }
    return validRecords
}

func GenerateReport(records []DataRecord) {
    fmt.Printf("Total records processed: %d\n", len(records))
    activeCount := 0
    for _, record := range records {
        if record.Active == "true" {
            activeCount++
        }
    }
    fmt.Printf("Active records: %d\n", activeCount)
    fmt.Println("--- Record Details ---")
    for i, record := range records {
        fmt.Printf("%d. ID: %s, Name: %s, Email: %s\n", i+1, record.ID, record.Name, record.Email)
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_processor.go <filename.csv>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := ProcessCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    validRecords := ValidateRecords(records)
    GenerateReport(validRecords)
}
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if len(strings.TrimSpace(data.Username)) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}
	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}
	return userData, nil
}