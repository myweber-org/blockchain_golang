
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

    timestamp, err := time.Parse("2006-01-02T15:04:05Z", parts[0])
    if err != nil {
        return LogEntry{}, err
    }

    return LogEntry{
        Timestamp: timestamp,
        Level:     parts[1],
        Message:   parts[2],
    }, nil
}

func filterLogsByLevel(entries []LogEntry, level string) []LogEntry {
    var filtered []LogEntry
    for _, entry := range entries {
        if entry.Level == level {
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
            entries = append(entries, entry)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return entries, nil
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

    fmt.Printf("Total log entries: %d\n", len(entries))

    errorLogs := filterLogsByLevel(entries, "ERROR")
    fmt.Printf("Error entries: %d\n", len(errorLogs))

    for _, entry := range errorLogs {
        fmt.Printf("[%s] %s: %s\n",
            entry.Timestamp.Format("2006-01-02 15:04:05"),
            entry.Level,
            entry.Message)
    }
}