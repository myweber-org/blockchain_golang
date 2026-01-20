package main

import (
	"fmt"
	"sort"
)

func CleanData(data []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range data {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	sort.Strings(result)
	return result
}

func main() {
	rawData := []string{"zebra", "apple", "banana", "apple", "cherry", "banana"}
	cleaned := CleanData(rawData)
	fmt.Println("Cleaned data:", cleaned)
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}