
package datautils

import "sort"

func CleanStringSlice(input []string) []string {
    if len(input) == 0 {
        return []string{}
    }

    seen := make(map[string]struct{})
    result := make([]string, 0, len(input))

    for _, item := range input {
        trimmed := strings.TrimSpace(item)
        if trimmed == "" {
            continue
        }
        if _, exists := seen[trimmed]; !exists {
            seen[trimmed] = struct{}{}
            result = append(result, trimmed)
        }
    }

    sort.Strings(result)
    return result
}package main

import (
	"fmt"
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func TrimAll(slice []string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = strings.TrimSpace(s)
	}
	return result
}

func CleanData(input []string) []string {
	trimmed := TrimAll(input)
	unique := RemoveDuplicates(trimmed)
	return unique
}

func main() {
	data := []string{"apple", " banana", "apple", " cherry ", "banana"}
	cleaned := CleanData(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}