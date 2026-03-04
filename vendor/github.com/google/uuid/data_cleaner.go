package utils

import "strings"

func TrimWhitespaceFromSlice(slice []string) []string {
    trimmed := make([]string, len(slice))
    for i, s := range slice {
        trimmed[i] = strings.TrimSpace(s)
    }
    return trimmed
}