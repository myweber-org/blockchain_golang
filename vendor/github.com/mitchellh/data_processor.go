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