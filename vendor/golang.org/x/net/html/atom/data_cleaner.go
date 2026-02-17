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
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package datautils

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
	if len(slice) < 2 {
		return slice
	}

	sort.Slice(slice, func(i, j int) bool {
		switch v := any(slice[i]).(type) {
		case int:
			return v < any(slice[j]).(int)
		case string:
			return v < any(slice[j]).(string)
		default:
			return false
		}
	})

	j := 0
	for i := 1; i < len(slice); i++ {
		if slice[j] != slice[i] {
			j++
			slice[j] = slice[i]
		}
	}

	return slice[:j+1]
}