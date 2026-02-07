package main

import (
	"errors"
	"sync"
)

type DataRecord struct {
	ID    int
	Value string
	Valid bool
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

func (p *Processor) AddRecord(id int, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.records = append(p.records, DataRecord{ID: id, Value: value, Valid: false})
}

func (p *Processor) ValidateRecord(id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := range p.records {
		if p.records[i].ID == id {
			if len(p.records[i].Value) > 0 {
				p.records[i].Valid = true
				return nil
			}
			return errors.New("empty value field")
		}
	}
	return errors.New("record not found")
}

func (p *Processor) ProcessAll() []DataRecord {
	p.mu.RLock()
	defer p.mu.RUnlock()

	validRecords := make([]DataRecord, 0)
	for _, record := range p.records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	proc := NewProcessor()
	proc.AddRecord(1, "alpha")
	proc.AddRecord(2, "beta")
	proc.AddRecord(3, "")

	proc.ValidateRecord(1)
	proc.ValidateRecord(2)
	proc.ValidateRecord(3)

	results := proc.ProcessAll()
	for _, r := range results {
		println(r.ID, r.Value)
	}
}
package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Valid     bool      `json:"valid"`
}

type Processor struct {
	mu      sync.RWMutex
	records map[string]DataRecord
}

func NewProcessor() *Processor {
	return &Processor{
		records: make(map[string]DataRecord),
	}
}

func (p *Processor) AddRecord(id string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	record := DataRecord{
		ID:        id,
		Value:     value,
		Timestamp: time.Now().UTC(),
		Valid:     value >= 0 && value <= 100,
	}

	p.records[id] = record
}

func (p *Processor) ValidateRecords() []DataRecord {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var validRecords []DataRecord
	for _, record := range p.records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func (p *Processor) ProcessBatch(ids []string, values []float64) {
	var wg sync.WaitGroup
	for i := range ids {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p.AddRecord(ids[idx], values[idx])
		}(i)
	}
	wg.Wait()
}

func (p *Processor) ExportJSON() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return json.MarshalIndent(p.records, "", "  ")
}

func main() {
	processor := NewProcessor()

	ids := []string{"A001", "B002", "C003", "D004"}
	values := []float64{42.5, 150.0, -10.2, 75.3}

	processor.ProcessBatch(ids, values)

	validRecords := processor.ValidateRecords()
	fmt.Printf("Valid records count: %d\n", len(validRecords))

	jsonData, err := processor.ExportJSON()
	if err != nil {
		fmt.Printf("Export error: %v\n", err)
		return
	}

	fmt.Println("Exported data:")
	fmt.Println(string(jsonData))
}