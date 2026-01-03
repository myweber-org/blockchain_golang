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

func parseLogLine(line string) (LogEntry, bool) {
    pattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.+)$`
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(line)

    if len(matches) != 4 {
        return LogEntry{}, false
    }

    return LogEntry{
        Timestamp: matches[1],
        Level:     strings.ToUpper(matches[2]),
        Message:   matches[3],
    }, true
}

func filterErrors(entries []LogEntry) []LogEntry {
    var errors []LogEntry
    for _, entry := range entries {
        if entry.Level == "ERROR" || entry.Level == "FATAL" {
            errors = append(errors, entry)
        }
    }
    return errors
}

func processLogFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    var entries []LogEntry
    scanner := bufio.NewScanner(file)
    lineNumber := 0

    for scanner.Scan() {
        lineNumber++
        line := scanner.Text()
        if entry, ok := parseLogLine(line); ok {
            entries = append(entries, entry)
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    errorEntries := filterErrors(entries)
    fmt.Printf("Total log entries: %d\n", len(entries))
    fmt.Printf("Error entries found: %d\n\n", len(errorEntries))

    for _, entry := range errorEntries {
        fmt.Printf("[%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Message)
    }

    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_parser <logfile>")
        os.Exit(1)
    }

    filename := os.Args[1]
    if err := processLogFile(filename); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }
}