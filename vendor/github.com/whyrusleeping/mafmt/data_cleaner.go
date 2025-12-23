
package csvutils

import (
	"encoding/csv"
	"io"
	"strings"
)

func CleanCSVData(input io.Reader, output io.Writer) error {
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := make([]string, 0, len(record))
		hasData := false

		for _, field := range record {
			trimmed := strings.TrimSpace(field)
			cleanedRecord = append(cleanedRecord, trimmed)
			if trimmed != "" {
				hasData = true
			}
		}

		if hasData {
			if err := writer.Write(cleanedRecord); err != nil {
				return err
			}
		}
	}

	return nil
}