package main

import (
	"fmt"
	"sort"
)

func CleanData(data []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, item := range data {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	sort.Strings(result)
	return result
}

func main() {
	rawData := []string{"zebra", "apple", "banana", "apple", "cherry", "banana"}
	cleanedData := CleanData(rawData)
	fmt.Println("Cleaned Data:", cleanedData)
}