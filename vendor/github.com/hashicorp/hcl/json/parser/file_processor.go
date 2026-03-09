package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type DataRecord struct {
	ID        string    `json:"id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	Processed bool      `json:"processed"`
}

type Processor struct {
	mu          sync.RWMutex
	records     map[string]DataRecord
	workerCount int
}

func NewProcessor(workers int) *Processor {
	return &Processor{
		records:     make(map[string]DataRecord),
		workerCount: workers,
	}
}

func (p *Processor) LoadData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var records []DataRecord
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	p.mu.Lock()
	for _, record := range records {
		p.records[record.ID] = record
	}
	p.mu.Unlock()

	log.Printf("Loaded %d records from %s", len(records), filename)
	return nil
}

func (p *Processor) ProcessRecord(id string) error {
	p.mu.RLock()
	record, exists := p.records[id]
	p.mu.RUnlock()

	if !exists {
		return fmt.Errorf("record with ID %s not found", id)
	}

	time.Sleep(50 * time.Millisecond)

	p.mu.Lock()
	record.Processed = true
	record.Value = record.Value * 1.1
	p.records[id] = record
	p.mu.Unlock()

	return nil
}

func (p *Processor) RunConcurrentProcessing() {
	var wg sync.WaitGroup
	ids := p.getAllIDs()

	workChan := make(chan string, len(ids))

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for id := range workChan {
				if err := p.ProcessRecord(id); err != nil {
					log.Printf("Worker %d: Error processing %s: %v", workerID, id, err)
				} else {
					log.Printf("Worker %d: Successfully processed %s", workerID, id)
				}
			}
		}(i)
	}

	for _, id := range ids {
		workChan <- id
	}
	close(workChan)

	wg.Wait()
	log.Println("All records processed")
}

func (p *Processor) getAllIDs() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	ids := make([]string, 0, len(p.records))
	for id := range p.records {
		ids = append(ids, id)
	}
	return ids
}

func (p *Processor) GenerateReport() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	processed := 0
	totalValue := 0.0

	for _, record := range p.records {
		if record.Processed {
			processed++
			totalValue += record.Value
		}
	}

	fmt.Printf("Processing Report:\n")
	fmt.Printf("Total records: %d\n", len(p.records))
	fmt.Printf("Processed records: %d\n", processed)
	fmt.Printf("Total value: %.2f\n", totalValue)
}

func main() {
	processor := NewProcessor(4)

	if err := processor.LoadData("data.json"); err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	processor.RunConcurrentProcessing()
	elapsed := time.Since(start)

	processor.GenerateReport()
	log.Printf("Processing completed in %v", elapsed)
}