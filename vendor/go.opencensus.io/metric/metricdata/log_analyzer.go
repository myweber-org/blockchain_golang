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

func analyzeLogs(filepath string) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    levelCount := make(map[string]int)
    var entries []LogEntry

    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            continue
        }
        entries = append(entries, entry)
        levelCount[entry.Level]++
    }

    fmt.Println("Log Analysis Summary:")
    fmt.Println("=====================")
    for level, count := range levelCount {
        fmt.Printf("%s: %d entries\n", level, count)
    }

    if len(entries) > 0 {
        fmt.Printf("\nFirst log: %v\n", entries[0].Timestamp)
        fmt.Printf("Last log:  %v\n", entries[len(entries)-1].Timestamp)
    }

    return scanner.Err()
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