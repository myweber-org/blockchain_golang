package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func CleanCSVRow(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		field = strings.TrimSpace(field)
		field = removeExtraSpaces(field)
		field = normalizeQuotes(field)
		cleaned[i] = field
	}
	return cleaned
}

func removeExtraSpaces(s string) string {
	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(s, " ")
}

func normalizeQuotes(s string) string {
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
		s = strings.ReplaceAll(s, `""`, `"`)
	}
	return s
}

func ProcessCSV(reader io.Reader, writer io.Writer) error {
	csvReader := csv.NewReader(reader)
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := CleanCSVRow(record)
		if err := csvWriter.Write(cleanedRecord); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fmt.Println("CSV data cleaner package loaded")
}