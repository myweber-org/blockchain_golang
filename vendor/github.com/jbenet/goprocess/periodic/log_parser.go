
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

    if matches == nil {
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

    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            fmt.Printf("Skipping invalid line: %v\n", err)
            continue
        }
        entries = append(entries, *entry)
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Total log entries: %d\n", len(entries))

    errorLogs := filterLogsByLevel(entries, "ERROR")
    fmt.Printf("Error entries: %d\n", len(errorLogs))

    for _, entry := range errorLogs {
        fmt.Printf("[%s] %s\n", entry.Timestamp.Format("15:04:05"), entry.Message)
    }
}