
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID      string
    Name    string
    Email   string
    Active  string
}

func cleanString(s string) string {
    return strings.TrimSpace(strings.ToLower(s))
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

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

    csvReader := csv.NewReader(inFile)
    csvWriter := csv.NewWriter(outFile)
    defer csvWriter.Flush()

    headers, err := csvReader.Read()
    if err != nil {
        return fmt.Errorf("failed to read headers: %w", err)
    }

    if err := csvWriter.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read record: %w", err)
        }

        if len(record) < 4 {
            continue
        }

        cleanedRecord := DataRecord{
            ID:     cleanString(record[0]),
            Name:   cleanString(record[1]),
            Email:  cleanString(record[2]),
            Active: cleanString(record[3]),
        }

        if cleanedRecord.ID == "" || cleanedRecord.Name == "" {
            continue
        }

        if !validateEmail(cleanedRecord.Email) {
            continue
        }

        if cleanedRecord.Active != "true" && cleanedRecord.Active != "false" {
            cleanedRecord.Active = "false"
        }

        outputRow := []string{
            cleanedRecord.ID,
            cleanedRecord.Name,
            cleanedRecord.Email,
            cleanedRecord.Active,
        }

        if err := csvWriter.Write(outputRow); err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run data_cleaner.go <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %s to %s\n", inputFile, outputFile)
}