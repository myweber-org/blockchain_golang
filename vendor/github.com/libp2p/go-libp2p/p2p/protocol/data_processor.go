
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