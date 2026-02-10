package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

type LogSummary struct {
	TotalEntries int
	ErrorCount   int
	WarnCount    int
	InfoCount    int
	UniqueErrors map[string]int
}

func parseLogLine(line string) (LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     matches[2],
		Message:   matches[3],
	}, nil
}

func analyzeLogs(filename string) (LogSummary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return LogSummary{}, err
	}
	defer file.Close()

	summary := LogSummary{
		UniqueErrors: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		summary.TotalEntries++
		switch strings.ToUpper(entry.Level) {
		case "ERROR":
			summary.ErrorCount++
			summary.UniqueErrors[entry.Message]++
		case "WARN":
			summary.WarnCount++
		case "INFO":
			summary.InfoCount++
		}
	}

	return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
	fmt.Printf("Log Analysis Summary:\n")
	fmt.Printf("Total Entries: %d\n", summary.TotalEntries)
	fmt.Printf("Errors: %d\n", summary.ErrorCount)
	fmt.Printf("Warnings: %d\n", summary.WarnCount)
	fmt.Printf("Info Messages: %d\n", summary.InfoCount)

	if len(summary.UniqueErrors) > 0 {
		fmt.Printf("\nUnique Errors:\n")
		for err, count := range summary.UniqueErrors {
			fmt.Printf("  %s (occurrences: %d)\n", err, count)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	summary, err := analyzeLogs(os.Args[1])
	if err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}

	printSummary(summary)
}