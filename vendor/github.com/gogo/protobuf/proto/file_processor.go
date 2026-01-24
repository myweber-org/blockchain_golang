
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type RecordProcessor interface {
	Process(record []string) error
}

type CSVHandler struct {
	processor RecordProcessor
}

func NewCSVHandler(p RecordProcessor) *CSVHandler {
	return &CSVHandler{processor: p}
}

func (h *CSVHandler) ProcessFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		if err := h.processor.Process(record); err != nil {
			return fmt.Errorf("processor error: %w", err)
		}
	}
	return nil
}

type PrintProcessor struct{}

func (p PrintProcessor) Process(record []string) error {
	fmt.Printf("Processed: %s\n", strings.Join(record, " | "))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <csv_file>")
		os.Exit(1)
	}

	processor := PrintProcessor{}
	handler := NewCSVHandler(processor)

	if err := handler.ProcessFile(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}