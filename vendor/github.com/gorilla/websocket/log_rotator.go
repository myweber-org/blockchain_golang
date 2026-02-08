package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    if err := os.MkdirAll(basePath, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    fullPath := filepath.Join(basePath, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to stat log file: %w", err)
    }

    return &LogRotator{
        currentFile: file,
        currentSize: info.Size(),
        basePath:    basePath,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, fmt.Errorf("rotation failed: %w", err)
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return fmt.Errorf("failed to close current log file: %w", err)
    }

    timestamp := time.Now().Format("20060102_150405")
    oldPath := filepath.Join(lr.basePath, logFileName)
    newPath := filepath.Join(lr.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))

    if err := os.Rename(oldPath, newPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    lr.currentFile = file
    lr.currentSize = 0

    go lr.cleanupOldFiles()

    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := filepath.Join(lr.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    sort.Sort(sort.Reverse(sort.StringSlice(matches)))

    for i, match := range matches {
        if i >= maxBackupFiles {
            os.Remove(match)
        }
    }
}

func (lr *LogRotator) parseBackupNumber(filename string) int {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return 0
    }
    num, _ := strconv.Atoi(parts[len(parts)-1])
    return num
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("./logs")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    writer := io.MultiWriter(os.Stdout, rotator)

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        writer.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}