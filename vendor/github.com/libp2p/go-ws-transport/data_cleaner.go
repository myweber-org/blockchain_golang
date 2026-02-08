package datautils

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
package main

import (
	"strings"
)

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

func TrimSpaces(input []string) []string {
	result := make([]string, len(input))
	for i, item := range input {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func CleanData(data []string) []string {
	trimmed := TrimSpaces(data)
	cleaned := RemoveDuplicates(trimmed)
	return cleaned
}