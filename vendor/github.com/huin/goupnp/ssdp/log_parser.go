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

    if matches == nil {
        return nil, fmt.Errorf("invalid log format")
    }

    return &LogEntry{
        Timestamp: matches[1],
        Level:     matches[2],
        Message:   matches[3],
    }, nil
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
    lineNumber := 0

    for scanner.Scan() {
        lineNumber++
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            fmt.Printf("Warning: Line %d - %v\n", lineNumber, err)
            continue
        }
        entries = append(entries, *entry)
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    errorEntries := filterErrors(entries)
    fmt.Printf("Found %d error entries:\n", len(errorEntries))
    for _, entry := range errorEntries {
        fmt.Printf("[%s] %s\n", entry.Timestamp, entry.Message)
    }
}