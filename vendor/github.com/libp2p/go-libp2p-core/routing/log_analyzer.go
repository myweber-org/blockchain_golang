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

type LogSummary struct {
    TotalEntries int
    ErrorCount   int
    WarnCount    int
    InfoCount    int
    StartTime    time.Time
    EndTime      time.Time
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

func analyzeLogs(filepath string) (LogSummary, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return LogSummary{}, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    summary := LogSummary{}
    firstEntry := true

    for scanner.Scan() {
        entry, err := parseLogLine(scanner.Text())
        if err != nil {
            continue
        }

        summary.TotalEntries++

        switch strings.ToUpper(entry.Level) {
        case "ERROR":
            summary.ErrorCount++
        case "WARN":
            summary.WarnCount++
        case "INFO":
            summary.InfoCount++
        }

        if firstEntry {
            summary.StartTime = entry.Timestamp
            summary.EndTime = entry.Timestamp
            firstEntry = false
        } else {
            if entry.Timestamp.Before(summary.StartTime) {
                summary.StartTime = entry.Timestamp
            }
            if entry.Timestamp.After(summary.EndTime) {
                summary.EndTime = entry.Timestamp
            }
        }
    }

    return summary, scanner.Err()
}

func printSummary(summary LogSummary) {
    fmt.Println("=== Log Analysis Summary ===")
    fmt.Printf("Total entries: %d\n", summary.TotalEntries)
    fmt.Printf("Error level: %d\n", summary.ErrorCount)
    fmt.Printf("Warning level: %d\n", summary.WarnCount)
    fmt.Printf("Info level: %d\n", summary.InfoCount)
    fmt.Printf("Time range: %s to %s\n",
        summary.StartTime.Format("2006-01-02 15:04:05"),
        summary.EndTime.Format("2006-01-02 15:04:05"))
    fmt.Printf("Duration: %v\n", summary.EndTime.Sub(summary.StartTime))
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_analyzer <logfile>")
        os.Exit(1)
    }

    summary, err := analyzeLogs(os.Args[1])
    if err != nil {
        fmt.Printf("Error analyzing logs: %v\n", err)
        os.Exit(1)
    }

    printSummary(summary)
}