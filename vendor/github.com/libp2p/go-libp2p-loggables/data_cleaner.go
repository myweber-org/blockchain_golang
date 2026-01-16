
package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{" apple", "banana ", "  apple  ", "banana", "", "cherry"}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange", "banana"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID        int
    Name      string
    Email     string
    Age       int
    Valid     bool
    Errors    []string
}

func parseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []Record{}
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %v", lineNumber, err)
        }

        if lineNumber == 0 {
            lineNumber++
            continue
        }

        record := parseLine(line, lineNumber)
        records = append(records, record)
        lineNumber++
    }

    return records, nil
}

func parseLine(fields []string, lineNum int) Record {
    record := Record{Valid: true}
    
    if len(fields) < 4 {
        record.Valid = false
        record.Errors = append(record.Errors, fmt.Sprintf("line %d: insufficient fields", lineNum))
        return record
    }

    id, err := strconv.Atoi(strings.TrimSpace(fields[0]))
    if err != nil {
        record.Valid = false
        record.Errors = append(record.Errors, fmt.Sprintf("line %d: invalid ID format", lineNum))
    } else {
        record.ID = id
    }

    record.Name = strings.TrimSpace(fields[1])
    if record.Name == "" {
        record.Valid = false
        record.Errors = append(record.Errors, fmt.Sprintf("line %d: name cannot be empty", lineNum))
    }

    record.Email = strings.TrimSpace(fields[2])
    if !strings.Contains(record.Email, "@") {
        record.Valid = false
        record.Errors = append(record.Errors, fmt.Sprintf("line %d: invalid email format", lineNum))
    }

    age, err := strconv.Atoi(strings.TrimSpace(fields[3]))
    if err != nil || age < 0 || age > 120 {
        record.Valid = false
        record.Errors = append(record.Errors, fmt.Sprintf("line %d: invalid age value", lineNum))
    } else {
        record.Age = age
    }

    return record
}

func validateRecords(records []Record) ([]Record, []Record) {
    valid := []Record{}
    invalid := []Record{}

    for _, record := range records {
        if record.Valid {
            valid = append(valid, record)
        } else {
            invalid = append(invalid, record)
        }
    }

    return valid, invalid
}

func generateReport(valid, invalid []Record) {
    fmt.Printf("Data Cleaning Report\n")
    fmt.Printf("====================\n")
    fmt.Printf("Total records processed: %d\n", len(valid)+len(invalid))
    fmt.Printf("Valid records: %d\n", len(valid))
    fmt.Printf("Invalid records: %d\n", len(invalid))

    if len(invalid) > 0 {
        fmt.Printf("\nInvalid Records Details:\n")
        for _, record := range invalid {
            fmt.Printf("  Record ID: %d\n", record.ID)
            for _, err := range record.Errors {
                fmt.Printf("    - %s\n", err)
            }
        }
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_cleaner.go <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := parseCSVFile(filename)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    valid, invalid := validateRecords(records)
    generateReport(valid, invalid)
}