package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, string, bool) {
	sanitizedName := dp.SanitizeString(name)
	sanitizedEmail := dp.SanitizeString(email)
	isValidEmail := dp.ValidateEmail(sanitizedEmail)
	return sanitizedName, sanitizedEmail, isValidEmail
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) NormalizeEmail(email string) (string, bool) {
	cleaned := dp.CleanString(email)
	lowerEmail := strings.ToLower(cleaned)
	return lowerEmail, dp.ValidateEmail(lowerEmail)
}

func (dp *DataProcessor) ProcessInputs(names []string, emails []string) ([]string, []string) {
	var cleanNames []string
	var validEmails []string

	for _, name := range names {
		cleanNames = append(cleanNames, dp.CleanString(name))
	}

	for _, email := range emails {
		if normalized, valid := dp.NormalizeEmail(email); valid {
			validEmails = append(validEmails, normalized)
		}
	}

	return cleanNames, validEmails
}