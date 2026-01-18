package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type LogSummary struct {
    TotalLines   int
    ErrorCount   int
    WarningCount int
    InfoCount    int
    UniqueErrors map[string]int
}

func NewLogSummary() *LogSummary {
    return &LogSummary{
        UniqueErrors: make(map[string]int),
    }
}

func analyzeLogFile(filename string) (*LogSummary, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    summary := NewLogSummary()
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        summary.TotalLines++
        line := scanner.Text()

        switch {
        case strings.Contains(line, "ERROR"):
            summary.ErrorCount++
            extractErrorPattern(line, summary.UniqueErrors)
        case strings.Contains(line, "WARNING"):
            summary.WarningCount++
        case strings.Contains(line, "INFO"):
            summary.InfoCount++
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return summary, nil
}

func extractErrorPattern(line string, errorMap map[string]int) {
    parts := strings.Split(line, "ERROR:")
    if len(parts) > 1 {
        errorMsg := strings.TrimSpace(parts[1])
        if len(errorMsg) > 50 {
            errorMsg = errorMsg[:50] + "..."
        }
        errorMap[errorMsg]++
    }
}

func printSummary(summary *LogSummary) {
    fmt.Println("=== Log Analysis Summary ===")
    fmt.Printf("Total lines processed: %d\n", summary.TotalLines)
    fmt.Printf("INFO entries: %d\n", summary.InfoCount)
    fmt.Printf("WARNING entries: %d\n", summary.WarningCount)
    fmt.Printf("ERROR entries: %d\n", summary.ErrorCount)

    if len(summary.UniqueErrors) > 0 {
        fmt.Println("\nUnique error patterns:")
        for pattern, count := range summary.UniqueErrors {
            fmt.Printf("  [%d occurrences] %s\n", count, pattern)
        }
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: log_analyzer <logfile>")
        os.Exit(1)
    }

    filename := os.Args[1]
    summary, err := analyzeLogFile(filename)
    if err != nil {
        fmt.Printf("Error analyzing file: %v\n", err)
        os.Exit(1)
    }

    printSummary(summary)
}