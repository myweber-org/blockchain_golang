package main

import "fmt"

func calculateAverage(numbers []float64) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    var sum float64
    for _, num := range numbers {
        sum += num
    }
    
    return sum / float64(len(numbers))
}

func main() {
    data := []float64{10.5, 20.3, 15.7, 8.9, 12.1}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
}