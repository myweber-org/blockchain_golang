
package main

import (
	"fmt"
)

// CalculateMovingAverage returns a slice containing the moving average of the input slice.
// The windowSize parameter defines the number of elements to include in each average calculation.
// If windowSize is greater than the length of data, an empty slice is returned.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
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
	sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	averages := CalculateMovingAverage(sampleData, 3)
	fmt.Println("Moving averages with window size 3:", averages)
}