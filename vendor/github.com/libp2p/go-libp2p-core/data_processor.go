
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func processCSVFile(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headerProcessed := false
	var header []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV record: %w", err)
		}

		if !headerProcessed {
			header = record
			headerProcessed = true
			if err := writer.Write(header); err != nil {
				return fmt.Errorf("error writing header: %w", err)
			}
			continue
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			cleanedField = strings.ToValidUTF8(cleanedField, "")
			if cleanedField == "" {
				cleanedField = "N/A"
			}
			cleanedRecord[i] = cleanedField
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := processCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %s -> %s\n", inputFile, outputFile)
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	IPAddress string `json:"ip_address"`
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

func sanitizeUsername(username string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return reg.ReplaceAllString(username, "")
}

func validateIP(ip string) bool {
	ipRegex := `^(\d{1,3}\.){3}\d{1,3}$`
	matched, _ := regexp.MatchString(ipRegex, ip)
	return matched
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	data.Email = normalizeEmail(data.Email)
	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = sanitizeUsername(data.Username)
	if len(data.Username) < 3 {
		return nil, fmt.Errorf("username too short")
	}

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("invalid age value: %d", data.Age)
	}

	if !validateIP(data.IPAddress) {
		return nil, fmt.Errorf("invalid IP address format: %s", data.IPAddress)
	}

	return &data, nil
}

func main() {
	jsonData := []byte(`{
		"email": "  TEST@Example.COM  ",
		"username": "user_123!@#",
		"age": 25,
		"ip_address": "192.168.1.1"
	}`)

	processedData, err := processUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed data: %+v\n", processedData)
}