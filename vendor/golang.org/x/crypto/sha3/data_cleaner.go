
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

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateRecord(rec Record) error {
    if rec.ID <= 0 {
        return fmt.Errorf("invalid ID: %d", rec.ID)
    }
    if len(rec.Name) == 0 {
        return fmt.Errorf("empty name for ID: %d", rec.ID)
    }
    if !strings.Contains(rec.Email, "@") {
        return fmt.Errorf("invalid email format: %s", rec.Email)
    }
    if rec.Score < 0 || rec.Score > 100 {
        return fmt.Errorf("score out of range: %.2f", rec.Score)
    }
    return nil
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
    lineNum := 0

    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: read error: %v", lineNum, err))
            continue
        }

        if len(row) != 5 {
            errors = append(errors, fmt.Errorf("line %d: expected 5 columns, got %d", lineNum, len(row)))
            continue
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: invalid ID: %v", lineNum, err))
            continue
        }

        name := strings.TrimSpace(row[1])
        email := cleanEmail(row[2])

        active, err := strconv.ParseBool(row[3])
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: invalid active flag: %v", lineNum, err))
            continue
        }

        score, err := strconv.ParseFloat(row[4], 64)
        if err != nil {
            errors = append(errors, fmt.Errorf("line %d: invalid score: %v", lineNum, err))
            continue
        }

        record := Record{
            ID:     id,
            Name:   name,
            Email:  email,
            Active: active,
            Score:  score,
        }

        if err := validateRecord(record); err != nil {
            errors = append(errors, fmt.Errorf("line %d: validation failed: %v", lineNum, err))
            continue
        }

        records = append(records, record)
    }

    return records, errors
}

func generateSummary(records []Record) {
    var totalScore float64
    activeCount := 0
    emailDomains := make(map[string]int)

    for _, rec := range records {
        totalScore += rec.Score
        if rec.Active {
            activeCount++
        }
        parts := strings.Split(rec.Email, "@")
        if len(parts) == 2 {
            emailDomains[parts[1]]++
        }
    }

    avgScore := totalScore / float64(len(records))
    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Average score: %.2f\n", avgScore)
    fmt.Printf("Active users: %d\n", activeCount)
    fmt.Println("Email domain distribution:")
    for domain, count := range emailDomains {
        fmt.Printf("  %s: %d\n", domain, count)
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run data_cleaner.go <csv_file>")
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
        fmt.Printf("\nSuccessfully processed %d records:\n", len(records))
        generateSummary(records)
    } else {
        fmt.Println("No valid records found in the file.")
    }
}