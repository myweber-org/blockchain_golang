
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]Record, 0)

    for line := 1; ; line++ {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", line, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d", line)
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found")
    }

    return records, nil
}

func CalculateStats(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, r := range records {
        sum += r.Value
    }
    average := sum / float64(len(records))

    var variance float64
    for _, r := range records {
        diff := r.Value - average
        variance += diff * diff
    }
    stdDev := variance / float64(len(records))

    return average, stdDev
}

func ValidateRecord(r Record) error {
    if r.ID <= 0 {
        return errors.New("ID must be positive")
    }
    if r.Name == "" {
        return errors.New("name cannot be empty")
    }
    if r.Value < 0 {
        return errors.New("value cannot be negative")
    }
    return nil
}
package main

import (
	"fmt"
	"strings"
	"unicode"
)

func NormalizeUsername(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	var normalized strings.Builder
	for _, r := range trimmed {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			normalized.WriteRune(unicode.ToLower(r))
		} else {
			return "", fmt.Errorf("invalid character in username: %c", r)
		}
	}

	result := normalized.String()
	if len(result) < 3 {
		return "", fmt.Errorf("username must be at least 3 characters")
	}
	if len(result) > 20 {
		return "", fmt.Errorf("username cannot exceed 20 characters")
	}

	return result, nil
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func main() {
	usernames := []string{"User_123", "  test  ", "ab", "invalid@name", "verylongusernameexceedinglimit"}
	for _, u := range usernames {
		normalized, err := NormalizeUsername(u)
		if err != nil {
			fmt.Printf("Error for '%s': %v\n", u, err)
		} else {
			fmt.Printf("Original: '%s' -> Normalized: '%s'\n", u, normalized)
		}
	}

	emails := []string{"test@example.com", "invalid", "user@", "@domain.com", "user@domain"}
	for _, e := range emails {
		fmt.Printf("Email '%s' valid: %v\n", e, ValidateEmail(e))
	}
}