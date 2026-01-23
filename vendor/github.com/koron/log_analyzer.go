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

func analyzeLogs(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	levelCount := make(map[string]int)
	var earliest, latest time.Time

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			fmt.Printf("Error parsing line %d: %v\n", lineNumber, err)
			continue
		}

		levelCount[entry.Level]++

		if earliest.IsZero() || entry.Timestamp.Before(earliest) {
			earliest = entry.Timestamp
		}
		if latest.IsZero() || entry.Timestamp.After(latest) {
			latest = entry.Timestamp
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println("Log Analysis Summary")
	fmt.Println("====================")
	fmt.Printf("Time range: %s to %s\n", earliest.Format(time.RFC3339), latest.Format(time.RFC3339))
	fmt.Println("Log level distribution:")
	for level, count := range levelCount {
		fmt.Printf("  %s: %d\n", level, count)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_analyzer <logfile>")
		os.Exit(1)
	}

	if err := analyzeLogs(os.Args[1]); err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}
}