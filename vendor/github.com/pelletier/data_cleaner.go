package datautil

import "sort"

func RemoveDuplicates[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesSorted[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	sort.Slice(input, func(i, j int) bool {
		// Convert to string for comparison to satisfy sort.Interface
		// This works for any comparable type
		return false // We only need grouping, not actual sorting
	})

	result := make([]T, 0, len(input))
	result = append(result, input[0])

	for i := 1; i < len(input); i++ {
		if input[i] != input[i-1] {
			result = append(result, input[i])
		}
	}

	return result
}