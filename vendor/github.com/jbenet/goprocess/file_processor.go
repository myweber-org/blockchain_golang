
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int
	Content   string
	Valid     bool
	Timestamp time.Time
}

type Processor struct {
	records []DataRecord
	mu      sync.RWMutex
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
	}
}

func (p *Processor) AddRecord(content string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        len(p.records) + 1,
		Content:   content,
		Valid:     len(content) > 0,
		Timestamp: time.Now(),
	}

	p.records = append(p.records, record)
}

func (p *Processor) ValidateRecords() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.records) == 0 {
		return errors.New("no records to validate")
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(p.records))

	for i := range p.records {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			record := &p.records[idx]
			
			if !record.Valid {
				errorChan <- fmt.Errorf("record %d failed validation", record.ID)
				return
			}
			
			if len(record.Content) > 100 {
				record.Valid = false
				errorChan <- fmt.Errorf("record %d content too long", record.ID)
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	var validationErrors []error
	for err := range errorChan {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed with %d errors", len(validationErrors))
	}

	return nil
}

func (p *Processor) GetStats() (int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	validCount := 0
	for _, record := range p.records {
		if record.Valid {
			validCount++
		}
	}

	return len(p.records), validCount
}

func main() {
	processor := NewProcessor()

	sampleData := []string{
		"Sample data record one",
		"",
		"Another valid data entry",
		"This record has extremely long content that should fail validation because it exceeds the maximum allowed length of 100 characters",
		"Short",
	}

	for _, data := range sampleData {
		processor.AddRecord(data)
	}

	total, valid := processor.GetStats()
	fmt.Printf("Initial stats: %d total records, %d valid\n", total, valid)

	if err := processor.ValidateRecords(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}

	total, valid = processor.GetStats()
	fmt.Printf("Final stats: %d total records, %d valid\n", total, valid)
}