
package data_processor

import (
	"regexp"
	"strings"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dc.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return strings.ToLower(normalized)
}

func (dc *DataCleaner) RemoveSpecialChars(input string, keepPattern string) string {
	if keepPattern == "" {
		keepPattern = `[^a-zA-Z0-9\s]`
	}
	regex := regexp.MustCompile(keepPattern)
	return regex.ReplaceAllString(input, "")
}

func (dc *DataCleaner) Tokenize(input string, delimiter string) []string {
	if delimiter == "" {
		delimiter = " "
	}
	normalized := dc.NormalizeString(input)
	return strings.Split(normalized, delimiter)
}