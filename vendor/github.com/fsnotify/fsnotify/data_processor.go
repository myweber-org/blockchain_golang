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
}