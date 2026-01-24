package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values.
// windowSize defines the number of elements to include in each average calculation.
// Returns a slice of averages or an error if the window size is invalid.
func CalculateMovingAverage(data []float64, windowSize int) ([]float64, error) {
	if windowSize <= 0 {
		return nil, fmt.Errorf("window size must be positive, got %d", windowSize)
	}
	if len(data) < windowSize {
		return nil, fmt.Errorf("data length %d is less than window size %d", len(data), windowSize)
	}

	var result []float64
	var sum float64

	// Calculate initial sum for the first window
	for i := 0; i < windowSize; i++ {
		sum += data[i]
	}
	result = append(result, sum/float64(windowSize))

	// Slide the window and update the sum
	for i := windowSize; i < len(data); i++ {
		sum = sum - data[i-windowSize] + data[i]
		result = append(result, sum/float64(windowSize))
	}

	return result, nil
}

func main() {
	sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3

	averages, err := CalculateMovingAverage(sampleData, window)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Moving averages (window=%d): %v\n", window, averages)
}package main

import (
	"encoding/csv"
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		if len(row) < 3 {
			return nil, fmt.Errorf("invalid row length: %d", len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %w", err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value format: %w", err)
		}

		record := DataRecord{
			ID:    id,
			Name:  row[1],
			Value: value,
		}
		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID: %d", record.ID)
		}

		if record.Name == "" {
			return fmt.Errorf("empty name for ID: %d", record.ID)
		}

		if record.Value < 0 {
			return fmt.Errorf("negative value for ID: %d", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID: %d", record.ID)
		}
		seenIDs[record.ID] = true
	}

	return nil
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