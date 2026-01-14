
package utils

import "strings"

func SanitizeInput(input string) string {
    trimmed := strings.TrimSpace(input)
    return strings.Join(strings.Fields(trimmed), " ")
}