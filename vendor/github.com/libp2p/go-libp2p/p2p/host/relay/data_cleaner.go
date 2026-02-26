
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

func deduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    var unique []DataRecord
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            unique = append(unique, record)
        }
    }
    return unique
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func validatePhone(phone string) bool {
    if len(phone) != 10 {
        return false
    }
    for _, ch := range phone {
        if ch < '0' || ch > '9' {
            return false
        }
    }
    return true
}

func cleanData(records []DataRecord) []DataRecord {
    var cleaned []DataRecord
    uniqueRecords := deduplicateRecords(records)
    
    for _, record := range uniqueRecords {
        if validateEmail(record.Email) && validatePhone(record.Phone) {
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    sampleData := []DataRecord{
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 2, Email: "invalid-email", Phone: "9876543210"},
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 3, Email: "user@domain.org", Phone: "5551234"},
    }
    
    cleaned := cleanData(sampleData)
    fmt.Printf("Original records: %d\n", len(sampleData))
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
    }
}
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
    ID    int
    Name  string
    Email string
    Score float64
}

func cleanData(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("line %d: %v", lineNumber, err)
        }

        if len(row) != 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(row[1])
        if name == "" {
            continue
        }

        email := strings.TrimSpace(row[2])
        if !strings.Contains(email, "@") {
            continue
        }

        score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
        if err != nil {
            continue
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Email: email,
            Score: score,
        })
    }

    return records, nil
}

func validateRecords(records []Record) []Record {
    var valid []Record
    seen := make(map[int]bool)

    for _, r := range records {
        if r.Score < 0 || r.Score > 100 {
            continue
        }
        if seen[r.ID] {
            continue
        }
        seen[r.ID] = true
        valid = append(valid, r)
    }

    return valid
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_cleaner <csv_file>")
        return
    }

    records, err := cleanData(os.Args[1])
    if err != nil {
        fmt.Printf("Error cleaning data: %v\n", err)
        return
    }

    validRecords := validateRecords(records)
    fmt.Printf("Processed %d records, %d valid records found\n", len(records), len(validRecords))

    for _, r := range validRecords {
        fmt.Printf("ID: %d, Name: %s, Email: %s, Score: %.2f\n", r.ID, r.Name, r.Email, r.Score)
    }
}