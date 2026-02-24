
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

func (dp *DataProcessor) CleanString(input string) string {
    cleaned := strings.TrimSpace(input)
    cleaned = strings.ToLower(cleaned)
    return cleaned
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
    return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) RemoveSpecialChars(input string) string {
    re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
    return re.ReplaceAllString(input, "")
}

func (dp *DataProcessor) ProcessUserData(name, email string) (string, bool) {
    cleanName := dp.CleanString(name)
    cleanName = dp.RemoveSpecialChars(cleanName)
    isValidEmail := dp.ValidateEmail(email)
    
    return cleanName, isValidEmail
}