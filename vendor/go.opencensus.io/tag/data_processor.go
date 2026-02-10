
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Tags      []string
	Valid     bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value cannot be negative")
	}
	if len(record.Tags) == 0 {
		return errors.New("record must have at least one tag")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Tags = normalizeTags(record.Tags)
	transformed.Value = roundValue(record.Value)
	transformed.Valid = true
	return transformed
}

func normalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var result []string
	
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !uniqueTags[normalized] {
			uniqueTags[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func roundValue(value float64) float64 {
	return float64(int(value*100)) / 100
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		
		processed = append(processed, TransformRecord(record))
	}
	
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Timestamp: time.Now(),
			Value:     123.4567,
			Tags:      []string{"Sensor", "TEMPERATURE", "sensor"},
		},
		{
			ID:        "rec002",
			Timestamp: time.Now().Add(-1 * time.Hour),
			Value:     98.765,
			Tags:      []string{"pressure", "  PRESSURE  "},
		},
	}
	
	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	
	fmt.Printf("Successfully processed %d records\n", len(processed))
	for _, rec := range processed {
		fmt.Printf("Record %s: value=%.2f, tags=%v, valid=%v\n", 
			rec.ID, rec.Value, rec.Tags, rec.Valid)
	}
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func ProcessCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if len(row) != 4 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        active, err := strconv.ParseBool(row[3])
        if err != nil {
            return nil, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:     id,
            Name:   name,
            Value:  value,
            Active: active,
        })
    }

    return records, nil
}

func CalculateStats(records []Record) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var max float64
    activeCount := 0

    for i, record := range records {
        sum += record.Value
        if i == 0 || record.Value > max {
            max = record.Value
        }
        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, max, activeCount
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    records, err := ProcessCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    avg, max, active := CalculateStats(records)
    fmt.Printf("Processed %d records\n", len(records))
    fmt.Printf("Average value: %.2f\n", avg)
    fmt.Printf("Maximum value: %.2f\n", max)
    fmt.Printf("Active records: %d\n", active)
}