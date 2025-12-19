
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
	Severity  string
	Message   string
}

func parseLogLine(line string) (*LogEntry, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05", parts[0])
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Timestamp: timestamp,
		Severity:  parts[1],
		Message:   parts[2],
	}, nil
}

func filterLogsBySeverity(entries []LogEntry, severity string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range entries {
		if strings.EqualFold(entry.Severity, severity) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func readLogFile(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			continue
		}
		entries = append(entries, *entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile> [severity]")
		os.Exit(1)
	}

	filename := os.Args[1]
	entries, err := readLogFile(filename)
	if err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		severity := os.Args[2]
		entries = filterLogsBySeverity(entries, severity)
		fmt.Printf("Showing %d entries with severity '%s':\n", len(entries), severity)
	} else {
		fmt.Printf("Showing all %d log entries:\n", len(entries))
	}

	for _, entry := range entries {
		fmt.Printf("[%s] %s: %s\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			strings.ToUpper(entry.Severity),
			entry.Message)
	}
}