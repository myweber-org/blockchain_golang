
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
}