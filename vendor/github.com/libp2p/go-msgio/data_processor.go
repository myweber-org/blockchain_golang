
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
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID    string
	Name  string
	Email string
	Score int
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("empty CSV file")
	}

	var data []DataRecord
	for i, row := range records {
		if i == 0 {
			continue
		}

		if len(row) < 4 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
		}

		if score, err := parseInt(row[3]); err == nil {
			record.Score = score
		}

		if isValidRecord(record) {
			data = append(data, record)
		}
	}

	return data, nil
}

func isValidRecord(record DataRecord) bool {
	return record.ID != "" && record.Name != "" && strings.Contains(record.Email, "@")
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &result)
	return result, err
}

func WriteProcessedData(records []DataRecord, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Email", "Score", "Status"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		status := "PASS"
		if record.Score < 60 {
			status = "FAIL"
		}

		row := []string{
			record.ID,
			record.Name,
			record.Email,
			fmt.Sprintf("%d", record.Score),
			status,
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func FilterRecords(records []DataRecord, minScore int) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Score >= minScore {
			filtered = append(filtered, record)
		}
	}
	return filtered
}