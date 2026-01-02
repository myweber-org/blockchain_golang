
package main

import (
	"fmt"
)

// FilterAndDouble filters even numbers from a slice and doubles their values.
func FilterAndDouble(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 == 0 {
			result = append(result, num*2)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := FilterAndDouble(input)
	fmt.Printf("Input: %v\n", input)
	fmt.Printf("Filtered and Doubled: %v\n", output)
}