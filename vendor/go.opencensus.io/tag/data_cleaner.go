
package main

import (
	"strings"
	"unicode"
)

func SanitizeCSVField(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if r == '"' {
			builder.WriteRune('"')
			builder.WriteRune('"')
		} else if r == ',' || r == '\n' || r == '\r' {
			builder.WriteRune(' ')
		} else if unicode.IsGraphic(r) && !unicode.IsControl(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}
	result := builder.String()
	if strings.ContainsAny(result, ",\"\n\r") {
		return "\"" + result + "\""
	}
	return result
}