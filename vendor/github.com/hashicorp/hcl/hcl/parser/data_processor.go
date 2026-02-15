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
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
	Valid bool    `json:"valid"`
}

func processCSVFile(inputPath string) ([]Record, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid row length at line %d", lineNumber)
		}

		age, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("invalid age at line %d: %w", lineNumber, err)
		}

		score, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid score at line %d: %w", lineNumber, err)
		}

		valid := row[3] == "true"

		records = append(records, Record{
			Name:  row[0],
			Age:   age,
			Score: score,
			Valid: valid,
		})
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

func saveToFile(content, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.json>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	records, err := processCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	jsonOutput, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	err = saveToFile(jsonOutput, outputFile)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d records to %s\n", len(records), outputFile)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var user UserData
	if err := json.Unmarshal(rawData, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonStr := `{"id": 101, "name": "Alice", "email": "alice@example.com"}`
	user, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Parsed user: %+v\n", user)
}