
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func cleanCSV(inputPath, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    seen := make(map[string]bool)
    headerWritten := false

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        for i := range record {
            record[i] = strings.TrimSpace(record[i])
        }

        key := strings.Join(record, "|")
        if seen[key] {
            continue
        }
        seen[key] = true

        if !headerWritten {
            headerWritten = true
        }

        if err := writer.Write(record); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    err := cleanCSV(os.Args[1], os.Args[2])
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
    seen := make(map[string]struct{})
    result := []string{}

    for _, item := range input {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}

func main() {
    sample := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
    cleaned := RemoveDuplicates(sample)
    fmt.Println("Original:", sample)
    fmt.Println("Cleaned:", cleaned)
}