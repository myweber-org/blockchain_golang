
package main

import (
	"fmt"
)

// FilterAndDouble filters out even numbers from a slice and doubles the remaining odd numbers.
func FilterAndDouble(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 != 0 {
			result = append(result, num*2)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	output := FilterAndDouble(input)
	fmt.Println("Processed slice:", output)
}