
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func processCSVFile(inputPath string, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headerProcessed := false
	rowCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV record: %w", err)
		}

		if !headerProcessed {
			headerProcessed = true
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing header: %w", err)
			}
			continue
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = strings.TrimSpace(field)
			if cleanedRecord[i] == "" {
				cleanedRecord[i] = "N/A"
			}
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
		rowCount++
	}

	fmt.Printf("Processed %d data rows successfully\n", rowCount)
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run data_processor.go <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CSV processing completed successfully")
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if strings.TrimSpace(data.Username) == "" {
		errors = append(errors, "username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		errors = append(errors, "invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		errors = append(errors, "age must be between 0 and 150")
	}

	return len(errors) == 0, errors
}

func TransformUsername(data *UserData) {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(data UserData) (UserData, error) {
	TransformUsername(&data)

	if valid, errs := ValidateUserData(data); !valid {
		return data, fmt.Errorf("validation failed: %v", strings.Join(errs, "; "))
	}

	return data, nil
}

func main() {
	sampleData := UserData{
		Username: "  TestUser  ",
		Email:    "user@example.com",
		Age:      25,
	}

	result, err := ProcessUserInput(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed data: %+v\n", result)
}package main

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
	Tags      []string `json:"tags"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func FilterInactiveUsers(users []UserProfile) []UserProfile {
	var activeUsers []UserProfile
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func TransformUserData(users []UserProfile) ([]map[string]interface{}, error) {
	var transformed []map[string]interface{}
	
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			return nil, fmt.Errorf("invalid email for user %d", user.ID)
		}
		
		transformedUser := map[string]interface{}{
			"user_id":   user.ID,
			"username":  NormalizeUsername(user.Username),
			"email":     strings.ToLower(user.Email),
			"age_group": categorizeAge(user.Age),
			"tag_count": len(user.Tags),
			"status":    "active",
		}
		
		if !user.Active {
			transformedUser["status"] = "inactive"
		}
		
		transformed = append(transformed, transformedUser)
	}
	
	return transformed, nil
}

func categorizeAge(age int) string {
	switch {
	case age < 18:
		return "minor"
	case age >= 18 && age <= 35:
		return "young_adult"
	case age > 35 && age <= 60:
		return "adult"
	default:
		return "senior"
	}
}

func ProcessUserJSON(jsonData string) (string, error) {
	var users []UserProfile
	
	err := json.Unmarshal([]byte(jsonData), &users)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}
	
	activeUsers := FilterInactiveUsers(users)
	
	transformedData, err := TransformUserData(activeUsers)
	if err != nil {
		return "", fmt.Errorf("transformation failed: %v", err)
	}
	
	resultJSON, err := json.MarshalIndent(transformedData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %v", err)
	}
	
	return string(resultJSON), nil
}