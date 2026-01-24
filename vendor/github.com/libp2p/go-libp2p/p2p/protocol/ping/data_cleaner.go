
package utils

import "strings"

func CleanStringSlice(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

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