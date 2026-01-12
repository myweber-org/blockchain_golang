package main

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
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
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
    ID    int
    Name  string
    Email string
    Score float64
}

func cleanCSV(inputPath, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("skipping line %d: %v\n", lineNum, err)
            continue
        }

        if len(row) != 4 {
            fmt.Printf("skipping line %d: invalid column count\n", lineNum)
            continue
        }

        record, err := validateRecord(row)
        if err != nil {
            fmt.Printf("skipping line %d: %v\n", lineNum, err)
            continue
        }

        cleanedRow := []string{
            strconv.Itoa(record.ID),
            strings.TrimSpace(record.Name),
            strings.ToLower(strings.TrimSpace(record.Email)),
            fmt.Sprintf("%.2f", record.Score),
        }

        if err := writer.Write(cleanedRow); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}

func validateRecord(row []string) (Record, error) {
    var record Record

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil || id <= 0 {
        return record, fmt.Errorf("invalid ID: %s", row[0])
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" || len(name) > 100 {
        return record, fmt.Errorf("invalid name: %s", row[1])
    }
    record.Name = name

    email := strings.TrimSpace(row[2])
    if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
        return record, fmt.Errorf("invalid email: %s", row[2])
    }
    record.Email = email

    score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
    if err != nil || score < 0 || score > 100 {
        return record, fmt.Errorf("invalid score: %s", row[3])
    }
    record.Score = score

    return record, nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputPath := os.Args[1]
    outputPath := os.Args[2]

    if err := cleanCSV(inputPath, outputPath); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}