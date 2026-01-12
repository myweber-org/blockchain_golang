
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
}