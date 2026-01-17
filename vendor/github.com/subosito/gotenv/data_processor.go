
package data_processor

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor(allowedPattern string) (*DataProcessor, error) {
	compiledPattern, err := regexp.Compile(allowedPattern)
	if err != nil {
		return nil, err
	}
	return &DataProcessor{allowedPattern: compiledPattern}, nil
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.allowedPattern.FindString(trimmed)
}

func (dp *DataProcessor) ValidateInput(input string) bool {
	return dp.allowedPattern.MatchString(input)
}

func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := dp.CleanInput(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}