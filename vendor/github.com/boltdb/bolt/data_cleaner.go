package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type Record struct {
    ID      string
    Name    string
    Email   string
    Status  string
}

func cleanString(s string) string {
    return strings.TrimSpace(strings.ToLower(s))
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSV(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNum := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNum++
        if lineNum == 1 {
            continue
        }

        if len(line) < 4 {
            continue
        }

        record := Record{
            ID:     cleanString(line[0]),
            Name:   cleanString(line[1]),
            Email:  cleanString(line[2]),
            Status: cleanString(line[3]),
        }

        if record.ID == "" || !validateEmail(record.Email) {
            continue
        }

        records = append(records, record)
    }

    return records, nil
}

func writeCSV(outputPath string, records []Record) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{"id", "name", "email", "status"}
    if err := writer.Write(header); err != nil {
        return err
    }

    for _, record := range records {
        row := []string{
            record.ID,
            record.Name,
            record.Email,
            record.Status,
        }
        if err := writer.Write(row); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    records, err := processCSV("input.csv")
    if err != nil {
        fmt.Printf("Error processing CSV: %v\n", err)
        return
    }

    fmt.Printf("Processed %d valid records\n", len(records))

    if err := writeCSV("output.csv", records); err != nil {
        fmt.Printf("Error writing CSV: %v\n", err)
        return
    }

    fmt.Println("Data cleaning completed successfully")
}