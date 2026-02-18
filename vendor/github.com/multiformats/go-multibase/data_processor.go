
package data_processor

import (
	"errors"
	"strings"
	"time"
)

type UserData struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	Active    bool
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}
	
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	
	if data.CreatedAt.After(time.Now()) {
		return errors.New("creation date cannot be in the future")
	}
	
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func GenerateUserReport(users []UserData) map[string]int {
	report := make(map[string]int)
	
	for _, user := range users {
		if user.Active {
			report["active"]++
		} else {
			report["inactive"]++
		}
		
		if strings.HasSuffix(user.Email, ".com") {
			report["dotcom_domains"]++
		}
	}
	
	report["total"] = len(users)
	return report
}

func FilterActiveUsers(users []UserData) []UserData {
	var activeUsers []UserData
	
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	
	return activeUsers
}

func CalculateAverageID(users []UserData) float64 {
	if len(users) == 0 {
		return 0
	}
	
	var sum int
	for _, user := range users {
		sum += user.ID
	}
	
	return float64(sum) / float64(len(users))
}package main

import (
    "encoding/csv"
    "errors"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID        int
    Name      string
    Value     float64
    Timestamp string
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

        if len(row) != 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(row[1])

        value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
        if err != nil {
            continue
        }

        timestamp := strings.TrimSpace(row[3])

        record := DataRecord{
            ID:        id,
            Name:      name,
            Value:     value,
            Timestamp: timestamp,
        }

        records = append(records, record)
    }

    return records, nil
}

func ValidateRecord(record DataRecord) error {
    if record.ID <= 0 {
        return errors.New("invalid ID")
    }

    if record.Name == "" {
        return errors.New("name cannot be empty")
    }

    if record.Value < 0 {
        return errors.New("value cannot be negative")
    }

    if record.Timestamp == "" {
        return errors.New("timestamp cannot be empty")
    }

    return nil
}

func CalculateAverage(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0
    }

    var sum float64
    for _, record := range records {
        sum += record.Value
    }

    return sum / float64(len(records))
}

func FilterByThreshold(records []DataRecord, threshold float64) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Value >= threshold {
            filtered = append(filtered, record)
        }
    }
    return filtered
}package main

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
	return strings.TrimSpace(strings.ToLower(username))
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return len(email) > 5 && len(email) < 255
}

func ValidateAge(age int) bool {
	return age >= 0 && age <= 120
}

func SanitizeInput(input string) string {
	var result strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Username = SanitizeInput(NormalizeUsername(data.Username))
	
	if !ValidateEmail(data.Email) {
		return data, fmt.Errorf("invalid email format")
	}
	
	if !ValidateAge(data.Age) {
		return data, fmt.Errorf("age out of valid range")
	}
	
	return data, nil
}

func main() {
	sampleData := UserData{
		Username: "  JohnDoe123  ",
		Email:    "john@example.com",
		Age:      25,
	}
	
	processed, err := ProcessUserData(sampleData)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	
	fmt.Printf("Processed data: %+v\n", processed)
}