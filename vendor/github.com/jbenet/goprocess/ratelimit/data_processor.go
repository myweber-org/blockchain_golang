
package main

import (
	"fmt"
)

// CalculateMovingAverage returns a slice containing the moving average of the input slice.
// The windowSize parameter defines the number of elements to average over.
// If windowSize is greater than the length of the data slice, an empty slice is returned.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return []float64{}
	}

	result := make([]float64, 0, len(data)-windowSize+1)
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

	return result
}

func main() {
	// Example usage
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Data: %v\n", data)
	fmt.Printf("Moving Average (window=%d): %v\n", window, averages)
}