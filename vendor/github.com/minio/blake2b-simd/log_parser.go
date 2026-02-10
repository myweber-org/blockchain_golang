
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

func parseLogLine(line string) (*LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid log format")
	}

	return &LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}, nil
}

func filterLogsByLevel(entries []LogEntry, level string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range entries {
		if strings.EqualFold(entry.Level, level) {
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
		if err == nil {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func displayLogEntries(entries []LogEntry) {
	for _, entry := range entries {
		fmt.Printf("[%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Message)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile> [level]")
		os.Exit(1)
	}

	filename := os.Args[1]
	entries, err := readLogFile(filename)
	if err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) == 3 {
		level := os.Args[2]
		entries = filterLogsByLevel(entries, level)
	}

	fmt.Printf("Found %d log entries\n", len(entries))
	displayLogEntries(entries)
}