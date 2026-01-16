
package main

import "fmt"

func FilterAndDouble(numbers []int, threshold int) []int {
    var result []int
    for _, num := range numbers {
        if num > threshold {
            result = append(result, num*2)
        }
    }
    return result
}

func main() {
    input := []int{1, 5, 10, 15, 20}
    filtered := FilterAndDouble(input, 8)
    fmt.Println("Original:", input)
    fmt.Println("Filtered and doubled:", filtered)
}