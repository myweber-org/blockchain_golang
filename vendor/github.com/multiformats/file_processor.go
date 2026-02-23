package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
)

type FileProcessor struct {
    inputDir  string
    outputDir string
    workers   int
}

func NewFileProcessor(input, output string, workers int) *FileProcessor {
    return &FileProcessor{
        inputDir:  input,
        outputDir: output,
        workers:   workers,
    }
}

func (fp *FileProcessor) ProcessFiles() error {
    files, err := ioutil.ReadDir(fp.inputDir)
    if err != nil {
        return fmt.Errorf("failed to read input directory: %w", err)
    }

    jobs := make(chan string, len(files))
    results := make(chan error, len(files))
    var wg sync.WaitGroup

    for w := 0; w < fp.workers; w++ {
        wg.Add(1)
        go fp.worker(jobs, results, &wg)
    }

    for _, file := range files {
        if !file.IsDir() {
            jobs <- file.Name()
        }
    }
    close(jobs)

    wg.Wait()
    close(results)

    for err := range results {
        if err != nil {
            return err
        }
    }

    return nil
}

func (fp *FileProcessor) worker(jobs <-chan string, results chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()

    for filename := range jobs {
        inputPath := filepath.Join(fp.inputDir, filename)
        outputPath := filepath.Join(fp.outputDir, filename)

        data, err := ioutil.ReadFile(inputPath)
        if err != nil {
            results <- fmt.Errorf("failed to read file %s: %w", filename, err)
            continue
        }

        processedData := processContent(data)

        if err := os.MkdirAll(fp.outputDir, 0755); err != nil {
            results <- fmt.Errorf("failed to create output directory: %w", err)
            continue
        }

        if err := ioutil.WriteFile(outputPath, processedData, 0644); err != nil {
            results <- fmt.Errorf("failed to write file %s: %w", filename, err)
            continue
        }

        results <- nil
    }
}

func processContent(data []byte) []byte {
    processed := make([]byte, len(data))
    for i, b := range data {
        processed[i] = b ^ 0xFF
    }
    return processed
}

func main() {
    processor := NewFileProcessor("./input", "./output", 4)
    if err := processor.ProcessFiles(); err != nil {
        fmt.Printf("Processing failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("File processing completed successfully")
}
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

func (p *Processor) AddRecord(record DataRecord) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.records[record.ID] = record
}

func (p *Processor) ProcessRecord(id string) error {
	p.mu.RLock()
	record, exists := p.records[id]
	p.mu.RUnlock()

	if !exists {
		return fmt.Errorf("record %s not found", id)
	}

	time.Sleep(50 * time.Millisecond)

	p.mu.Lock()
	record.Processed = true
	p.records[id] = record
	p.mu.Unlock()

	return nil
}

func (p *Processor) ProcessAll() []error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(p.records))
	semaphore := make(chan struct{}, p.workerCount)

	p.mu.RLock()
	ids := make([]string, 0, len(p.records))
	for id := range p.records {
		ids = append(ids, id)
	}
	p.mu.RUnlock()

	for _, id := range ids {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(recordID string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if err := p.ProcessRecord(recordID); err != nil {
				errChan <- err
			}
		}(id)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	return errors
}

func (p *Processor) ExportToFile(filename string) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(p.records)
}

func generateSampleData(count int) []DataRecord {
	records := make([]DataRecord, count)
	for i := 0; i < count; i++ {
		records[i] = DataRecord{
			ID:        fmt.Sprintf("REC-%04d", i+1),
			Value:     float64(i) * 1.5,
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
			Processed: false,
		}
	}
	return records
}

func main() {
	logger := log.New(os.Stdout, "PROCESSOR: ", log.LstdFlags)
	processor := NewProcessor(4)

	sampleData := generateSampleData(100)
	for _, record := range sampleData {
		processor.AddRecord(record)
	}

	logger.Printf("Processing %d records with %d workers", len(sampleData), processor.workerCount)
	startTime := time.Now()

	errors := processor.ProcessAll()
	elapsed := time.Since(startTime)

	if len(errors) > 0 {
		logger.Printf("Completed with %d errors in %v", len(errors), elapsed)
		for _, err := range errors {
			logger.Printf("Error: %v", err)
		}
	} else {
		logger.Printf("All records processed successfully in %v", elapsed)
	}

	if err := processor.ExportToFile("processed_data.json"); err != nil {
		logger.Printf("Export failed: %v", err)
	} else {
		logger.Println("Data exported to processed_data.json")
	}
}