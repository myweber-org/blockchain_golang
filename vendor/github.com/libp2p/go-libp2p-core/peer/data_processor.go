
package main

import "fmt"

func FilterAndTransform(nums []int, threshold int) []int {
    var result []int
    for _, num := range nums {
        if num > threshold {
            transformed := num * 2
            result = append(result, transformed)
        }
    }
    return result
}

func main() {
    input := []int{1, 5, 10, 15, 20}
    output := FilterAndTransform(input, 8)
    fmt.Println("Processed slice:", output)
}