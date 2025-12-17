
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
    ID        int     `json:"id"`
    Name      string  `json:"name"`
    Value     float64 `json:"value"`
    Processed bool    `json:"processed"`
}

func processCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []Record
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 0 {
            lineNumber++
            continue
        }

        if len(line) < 3 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        record := Record{
            ID:        id,
            Name:      line[1],
            Value:     value,
            Processed: value > 50.0,
        }

        records = append(records, record)
        lineNumber++
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
    if len(os.Args) < 3 {
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
}