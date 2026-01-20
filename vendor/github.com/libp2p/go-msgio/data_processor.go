
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return &DataProcessor{emailRegex: regex}
}

func (dp *DataProcessor) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
	sanitizedName := dp.SanitizeInput(name)
	sanitizedEmail := dp.SanitizeInput(email)

	if sanitizedName == "" || !dp.ValidateEmail(sanitizedEmail) {
		return "", false
	}

	return sanitizedName + " <" + sanitizedEmail + ">", true
}