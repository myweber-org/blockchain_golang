
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID      int
    Name    string
    Email   string
    Active  bool
    Score   float64
}

func cleanString(s string) string {
    return strings.TrimSpace(s)
}

func parseBool(s string) (bool, error) {
    s = strings.ToLower(cleanString(s))
    return s == "true" || s == "yes" || s == "1" || s == "y", nil
}

func parseRecord(row []string) (Record, error) {
    if len(row) != 5 {
        return Record{}, fmt.Errorf("invalid row length: %d", len(row))
    }

    id, err := strconv.Atoi(cleanString(row[0]))
    if err != nil {
        return Record{}, fmt.Errorf("invalid ID: %v", err)
    }

    name := cleanString(row[1])
    if name == "" {
        return Record{}, fmt.Errorf("name cannot be empty")
    }

    email := cleanString(row[2])
    if !strings.Contains(email, "@") {
        return Record{}, fmt.Errorf("invalid email format")
    }

    active, err := parseBool(row[3])
    if err != nil {
        return Record{}, fmt.Errorf("invalid active flag: %v", err)
    }

    score, err := strconv.ParseFloat(cleanString(row[4]), 64)
    if err != nil {
        return Record{}, fmt.Errorf("invalid score: %v", err)
    }

    return Record{
        ID:     id,
        Name:   name,
        Email:  email,
        Active: active,
        Score:  score,
    }, nil
}

func processCSVFile(filename string) ([]Record, []error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, []error{err}
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []Record
    var errors []error
    line := 0

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        line++

        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: read error: %v", line, err))
            continue
        }

        if line == 1 {
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: %v", line, err))
            continue
        }

        records = append(records, record)
    }

    return records, errors
}

func generateSummary(records []Record) {
    var totalScore float64
    activeCount := 0
    uniqueDomains := make(map[string]bool)

    for _, r := range records {
        totalScore += r.Score
        if r.Active {
            activeCount++
        }

        parts := strings.Split(r.Email, "@")
        if len(parts) == 2 {
            uniqueDomains[parts[1]] = true
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Active records: %d\n", activeCount)
    if len(records) > 0 {
        fmt.Printf("Average score: %.2f\n", totalScore/float64(len(records)))
    }
    fmt.Printf("Unique email domains: %d\n", len(uniqueDomains))
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: data_cleaner <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, errors := processCSVFile(filename)

    if len(errors) > 0 {
        fmt.Printf("Encountered %d errors during processing:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("  - %v\n", err)
        }
    }

    if len(records) > 0 {
        fmt.Println("\nValid records:")
        for _, r := range records {
            fmt.Printf("  ID: %d, Name: %s, Email: %s, Active: %v, Score: %.1f\n",
                r.ID, r.Name, r.Email, r.Active, r.Score)
        }

        fmt.Println("\nSummary statistics:")
        generateSummary(records)
    } else {
        fmt.Println("No valid records found")
    }
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)
}