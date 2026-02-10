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

func parseLogLine(line string) *LogEntry {
    re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`)
    matches := re.FindStringSubmatch(line)
    
    if matches == nil {
        return nil
    }
    
    return &LogEntry{
        Timestamp: matches[1],
        Level:     matches[2],
        Message:   matches[3],
    }
}

func filterErrors(entries []LogEntry) []LogEntry {
    var errorEntries []LogEntry
    for _, entry := range entries {
        if strings.ToUpper(entry.Level) == "ERROR" {
            errorEntries = append(errorEntries, entry)
        }
    }
    return errorEntries
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
        if entry := parseLogLine(scanner.Text()); entry != nil {
            entries = append(entries, *entry)
        }
    }
    
    return entries, scanner.Err()
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_parser <logfile>")
        os.Exit(1)
    }

    entries, err := readLogFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error reading log file: %v\n", err)
        os.Exit(1)
    }

    errorEntries := filterErrors(entries)
    
    fmt.Printf("Total log entries: %d\n", len(entries))
    fmt.Printf("Error entries: %d\n\n", len(errorEntries))
    
    for _, entry := range errorEntries {
        fmt.Printf("[%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Message)
    }
}
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

func parseLogLine(line string) (*LogEntry, error) {
	pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid log format")
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Timestamp: timestamp,
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
		if err != nil {
			continue
		}
		entries = append(entries, *entry)
	}

	return entries, scanner.Err()
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

	if len(os.Args) > 2 {
		level := os.Args[2]
		entries = filterLogsByLevel(entries, level)
	}

	for _, entry := range entries {
		fmt.Printf("%s [%s] %s\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Level,
			entry.Message)
	}
}