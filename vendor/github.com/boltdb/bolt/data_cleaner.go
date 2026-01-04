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
}package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID   int
	Name string
	Age  int
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	result := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%d-%s-%d", record.ID, strings.ToLower(record.Name), record.Age)
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateRecords() (valid []DataRecord, invalid []DataRecord) {
	valid = make([]DataRecord, 0)
	invalid = make([]DataRecord, 0)

	for _, record := range dc.records {
		if record.ID > 0 && record.Name != "" && record.Age >= 0 && record.Age <= 120 {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}

	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(DataRecord{ID: 1, Name: "John", Age: 30})
	cleaner.AddRecord(DataRecord{ID: 2, Name: "Jane", Age: 25})
	cleaner.AddRecord(DataRecord{ID: 1, Name: "John", Age: 30})
	cleaner.AddRecord(DataRecord{ID: 3, Name: "", Age: 40})
	cleaner.AddRecord(DataRecord{ID: 4, Name: "Bob", Age: 150})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.GetRecordCount())

	valid, invalid := cleaner.ValidateRecords()
	fmt.Printf("Valid records: %d\n", len(valid))
	fmt.Printf("Invalid records: %d\n", len(invalid))

	for _, record := range invalid {
		fmt.Printf("Invalid: ID=%d, Name='%s', Age=%d\n", record.ID, record.Name, record.Age)
	}
}