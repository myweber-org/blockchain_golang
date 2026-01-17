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
}