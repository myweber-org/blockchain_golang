package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, value := range input {
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}

func main() {
	slice := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	uniqueSlice := RemoveDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Unique:", uniqueSlice)
}