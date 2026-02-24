
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

type Processor struct {
	validationRules map[string]func(float64) bool
	transformations []func(*DataRecord) error
}

func NewProcessor() *Processor {
	return &Processor{
		validationRules: make(map[string]func(float64) bool),
		transformations: make([]func(*DataRecord) error, 0),
	}
}

func (p *Processor) AddValidation(name string, rule func(float64) bool) {
	p.validationRules[name] = rule
}

func (p *Processor) AddTransformation(t func(*DataRecord) error) {
	p.transformations = append(p.transformations, t)
}

func (p *Processor) ValidateValue(value float64) error {
	for name, rule := range p.validationRules {
		if !rule(value) {
			return fmt.Errorf("validation failed: %s", name)
		}
	}
	return nil
}

func (p *Processor) ProcessRecord(record *DataRecord) error {
	if err := p.ValidateValue(record.Value); err != nil {
		return err
	}

	for _, transform := range p.transformations {
		if err := transform(record); err != nil {
			return fmt.Errorf("transformation error: %w", err)
		}
	}

	return nil
}

func NormalizeTags(record *DataRecord) error {
	uniqueTags := make(map[string]bool)
	var normalized []string

	for _, tag := range record.Tags {
		trimmed := strings.TrimSpace(strings.ToLower(tag))
		if trimmed != "" && !uniqueTags[trimmed] {
			uniqueTags[trimmed] = true
			normalized = append(normalized, trimmed)
		}
	}

	record.Tags = normalized
	return nil
}

func ApplyLogTransform(record *DataRecord) error {
	if record.Value <= 0 {
		return errors.New("value must be positive for log transformation")
	}
	record.Value = record.Value * 0.5
	return nil
}

func main() {
	processor := NewProcessor()

	processor.AddValidation("positive", func(v float64) bool { return v > 0 })
	processor.AddValidation("reasonable_range", func(v float64) bool { return v < 10000 })

	processor.AddTransformation(NormalizeTags)
	processor.AddTransformation(ApplyLogTransform)

	sampleRecord := &DataRecord{
		ID:        "rec-001",
		Value:     42.5,
		Timestamp: time.Now(),
		Tags:      []string{"  TEMP  ", "sensor", "TEMP", "data"},
	}

	fmt.Printf("Original record: %+v\n", sampleRecord)

	if err := processor.ProcessRecord(sampleRecord); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	fmt.Printf("Processed record: %+v\n", sampleRecord)
	fmt.Printf("Tags after normalization: %v\n", sampleRecord.Tags)
}
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSVFile(filepath string) ([]DataRecord, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
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

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	if len(records) == 0 {
		return nil, errors.New("no valid records found in CSV file")
	}

	return records, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func FilterByThreshold(records []DataRecord, threshold float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value >= threshold {
			filtered = append(filtered, record)
		}
	}
	return filtered
}