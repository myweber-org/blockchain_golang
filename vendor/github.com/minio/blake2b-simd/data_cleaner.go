package main

import (
	"regexp"
	"strings"
)

func SanitizeCSVField(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove any non-printable characters except comma, period, dash, and underscore
	reg := regexp.MustCompile(`[^\x20-\x7E,.\-_]`)
	cleaned := reg.ReplaceAllString(trimmed, "")
	
	// Replace multiple spaces with single space
	spaceReg := regexp.MustCompile(`\s+`)
	final := spaceReg.ReplaceAllString(cleaned, " ")
	
	return final
}