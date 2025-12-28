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