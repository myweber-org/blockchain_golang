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

func processCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record

    // Skip header
    _, err = reader.Read()
    if err != nil {
        return nil, err
    }

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) < 4 {
            continue
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            continue
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            continue
        }

        processed, err := strconv.ParseBool(row[3])
        if err != nil {
            processed = false
        }

        record := Record{
            ID:        id,
            Name:      row[1],
            Value:     value,
            Processed: processed,
        }
        records = append(records, record)
    }

    return records, nil
}

func convertToJSON(records []Record) (string, error) {
    jsonData, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
}

func filterProcessed(records []Record) []Record {
    var filtered []Record
    for _, record := range records {
        if record.Processed {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func calculateTotal(records []Record) float64 {
    var total float64
    for _, record := range records {
        total += record.Value
    }
    return total
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        return
    }

    records, err := processCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    fmt.Printf("Total records: %d\n", len(records))
    fmt.Printf("Total value: %.2f\n", calculateTotal(records))

    processedRecords := filterProcessed(records)
    fmt.Printf("Processed records: %d\n", len(processedRecords))

    jsonOutput, err := convertToJSON(processedRecords)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        return
    }

    fmt.Println("\nJSON Output:")
    fmt.Println(jsonOutput)
}