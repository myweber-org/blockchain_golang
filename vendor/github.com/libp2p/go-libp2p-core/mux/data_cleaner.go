
package main

import "fmt"

func removeDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    return result
}

func main() {
    input := []int{1, 2, 2, 3, 4, 4, 5, 6, 6, 7}
    cleaned := removeDuplicates(input)
    fmt.Printf("Original: %v\n", input)
    fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

func CleanStringSlice(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{"  apple ", "banana", "  apple", "banana ", "  ", "cherry"}
	cleaned := CleanStringSlice(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}