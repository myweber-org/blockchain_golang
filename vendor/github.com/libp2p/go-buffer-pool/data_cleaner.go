package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataCleaner struct {
    inputPath  string
    outputPath string
    delimiter  rune
}

func NewDataCleaner(input, output string) *DataCleaner {
    return &DataCleaner{
        inputPath:  input,
        outputPath: output,
        delimiter:  ',',
    }
}

func (dc *DataCleaner) SetDelimiter(delim rune) {
    dc.delimiter = delim
}

func (dc *DataCleaner) RemoveDuplicates() error {
    inFile, err := os.Open(dc.inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(dc.outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    reader.Comma = dc.delimiter
    writer := csv.NewWriter(outFile)
    writer.Comma = dc.delimiter
    defer writer.Flush()

    seen := make(map[string]bool)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read CSV record: %w", err)
        }

        key := strings.Join(record, "|")
        if !seen[key] {
            seen[key] = true
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("failed to write CSV record: %w", err)
            }
        }
    }

    return nil
}

func main() {
    cleaner := NewDataCleaner("input.csv", "output.csv")
    cleaner.SetDelimiter(',')
    
    if err := cleaner.RemoveDuplicates(); err != nil {
        fmt.Printf("Error cleaning data: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Data cleaning completed successfully")
}
package main

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
	data := []int{4, 2, 8, 2, 4, 9, 8, 1}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}