package main

import (
	"bufio"
	"fmt"
	"os"
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
	LevelCounts  map[string]int
	Errors       []string
}

func parseLogLine(line string) (LogEntry, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return LogEntry{}, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05", parts[0])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     parts[1],
		Message:   parts[2],
	}, nil
}

func analyzeLogFile(filename string) (LogSummary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return LogSummary{}, err
	}
	defer file.Close()

	summary := LogSummary{
		LevelCounts: make(map[string]int),
		Errors:      []string{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}

		summary.TotalEntries++
		summary.LevelCounts[entry.Level]++

		if entry.Level == "ERROR" {
			summary.Errors = append(summary.Errors, entry.Message)
		}
	}

	return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
	fmt.Printf("Log Analysis Summary:\n")
	fmt.Printf("Total entries: %d\n", summary.TotalEntries)
	fmt.Printf("Level distribution:\n")
	for level, count := range summary.LevelCounts {
		fmt.Printf("  %s: %d\n", level, count)
	}
	if len(summary.Errors) > 0 {
		fmt.Printf("Error messages found:\n")
		for _, err := range summary.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	summary, err := analyzeLogFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error analyzing log file: %v\n", err)
		os.Exit(1)
	}

	printSummary(summary)
}