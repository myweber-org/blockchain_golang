package main

import "fmt"

func calculateMovingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || windowSize > len(data) {
        return nil
    }

    result := make([]float64, 0, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := i; j < i+windowSize; j++ {
            sum += data[j]
        }
        result = append(result, sum/float64(windowSize))
    }
    return result
}

func main() {
    testData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := calculateMovingAverage(testData, window)
    fmt.Printf("Moving averages with window size %d: %v\n", window, averages)
}