
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type FileRecord struct {
	ID   int
	Data string
}

func validateRecord(record FileRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid record ID")
	}
	if len(record.Data) == 0 {
		return errors.New("empty data field")
	}
	return nil
}

func processFile(filename string) ([]FileRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var records []FileRecord
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		records = append(records, FileRecord{
			ID:   lineNumber,
			Data: scanner.Text(),
		})
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return records, nil
}

func concurrentValidation(records []FileRecord) ([]FileRecord, []error) {
	var validRecords []FileRecord
	var validationErrors []error
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, record := range records {
		wg.Add(1)
		go func(r FileRecord) {
			defer wg.Done()
			err := validateRecord(r)
			
			mu.Lock()
			if err != nil {
				validationErrors = append(validationErrors, 
					fmt.Errorf("record %d: %w", r.ID, err))
			} else {
				validRecords = append(validRecords, r)
			}
			mu.Unlock()
		}(record)
	}

	wg.Wait()
	return validRecords, validationErrors
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <filename>")
		os.Exit(1)
	}

	startTime := time.Now()
	filename := os.Args[1]

	records, err := processFile(filename)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d records from %s\n", len(records), filename)

	validRecords, errors := concurrentValidation(records)
	
	fmt.Printf("Validation completed in %v\n", time.Since(startTime))
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Validation errors: %d\n", len(errors))

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("Error: %v\n", err)
		}
	}
}package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type DataRecord struct {
	ID        int       `json:"id"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Processed bool      `json:"processed"`
}

type Processor struct {
	mu       sync.RWMutex
	records  []DataRecord
	errors   []error
	wg       sync.WaitGroup
}

func NewProcessor() *Processor {
	return &Processor{
		records: make([]DataRecord, 0),
		errors:  make([]error, 0),
	}
}

func (p *Processor) AddRecord(record DataRecord) {
	p.mu.Lock()
	defer p.mu.Unlock()
	record.Timestamp = time.Now()
	p.records = append(p.records, record)
}

func (p *Processor) ProcessRecord(index int) {
	defer p.wg.Done()

	p.mu.RLock()
	if index >= len(p.records) {
		p.mu.RUnlock()
		return
	}
	record := p.records[index]
	p.mu.RUnlock()

	time.Sleep(10 * time.Millisecond)

	p.mu.Lock()
	p.records[index].Processed = true
	p.mu.Unlock()
}

func (p *Processor) ProcessAll() {
	p.wg.Add(len(p.records))
	for i := range p.records {
		go p.ProcessRecord(i)
	}
	p.wg.Wait()
}

func (p *Processor) LogError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.errors = append(p.errors, err)
	log.Printf("Error recorded: %v", err)
}

func (p *Processor) ExportToFile(filename string) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		p.LogError(err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(p.records); err != nil {
		p.LogError(err)
		return err
	}
	return nil
}

func (p *Processor) Stats() (int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	processed := 0
	for _, record := range p.records {
		if record.Processed {
			processed++
		}
	}
	return len(p.records), processed
}

func main() {
	processor := NewProcessor()

	for i := 1; i <= 100; i++ {
		record := DataRecord{
			ID:    i,
			Value: fmt.Sprintf("data-%d", i),
		}
		processor.AddRecord(record)
	}

	processor.ProcessAll()

	total, processed := processor.Stats()
	fmt.Printf("Processed %d out of %d records\n", processed, total)

	if err := processor.ExportToFile("output.json"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data exported to output.json")
}