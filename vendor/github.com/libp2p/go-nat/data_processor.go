
package data_processor

import (
	"regexp"
	"strings"
	"unicode"
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
	return normalized
}

func (dc *DataCleaner) RemoveSpecialCharacters(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, input)
}

func (dc *DataCleaner) Process(input string) string {
	cleaned := dc.RemoveSpecialCharacters(input)
	return dc.NormalizeString(cleaned)
}