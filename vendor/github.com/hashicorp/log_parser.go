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
		Level:     strings.ToUpper(matches[2]),
		Message:   matches[3],
	}, nil
}

func filterErrors(entries []LogEntry) []LogEntry {
	var errorEntries []LogEntry
	for _, entry := range entries {
		if entry.Level == "ERROR" {
			errorEntries = append(errorEntries, entry)
		}
	}
	return errorEntries
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log_parser <logfile>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			fmt.Printf("Warning: Line %d - %v\n", lineNumber, err)
		} else {
			entries = append(entries, *entry)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	errorEntries := filterErrors(entries)
	fmt.Printf("Total log entries: %d\n", len(entries))
	fmt.Printf("Error entries found: %d\n\n", len(errorEntries))

	for i, entry := range errorEntries {
		fmt.Printf("Error #%d\n", i+1)
		fmt.Printf("Time: %s\n", entry.Timestamp)
		fmt.Printf("Message: %s\n\n", entry.Message)
	}
}