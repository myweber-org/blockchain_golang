
package main

import (
	"fmt"
	"sort"
)

func FilterAndSort(numbers []int, threshold int) []int {
	var filtered []int
	for _, num := range numbers {
		if num > threshold {
			filtered = append(filtered, num)
		}
	}
	sort.Ints(filtered)
	return filtered
}

func main() {
	data := []int{45, 12, 89, 3, 67, 34, 9, 21}
	result := FilterAndSort(data, 20)
	fmt.Println("Filtered and sorted:", result)
}