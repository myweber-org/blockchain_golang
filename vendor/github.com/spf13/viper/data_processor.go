
package main

import (
	"fmt"
)

// MovingAverage calculates the moving average of a slice of float64 numbers.
// It returns a new slice where each element is the average of the previous 'windowSize' elements.
// If the slice length is less than windowSize, it returns an empty slice.
func MovingAverage(data []float64, windowSize int) []float64 {
	if len(data) < windowSize || windowSize <= 0 {
		return []float64{}
	}

	result := make([]float64, len(data)-windowSize+1)
	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += data[i+j]
		}
		result[i] = sum / float64(windowSize)
	}
	return result
}

func main() {
	// Example usage
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := MovingAverage(data, window)
	fmt.Printf("Moving averages with window size %d: %v\n", window, averages)
}