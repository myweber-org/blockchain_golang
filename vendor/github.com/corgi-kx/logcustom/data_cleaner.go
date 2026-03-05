package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package datautils

import "sort"

func RemoveDuplicates[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesSorted[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	sort.Slice(slice, func(i, j int) bool {
		switch v := any(slice).(type) {
		case []string:
			return v[i] < v[j]
		case []int:
			return v[i] < v[j]
		case []float64:
			return v[i] < v[j]
		default:
			return false
		}
	})

	result := make([]T, 0, len(slice))
	result = append(result, slice[0])

	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1] {
			result = append(result, slice[i])
		}
	}

	return result
}