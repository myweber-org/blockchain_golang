package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
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

func NewLogRotator(logDir string) (*LogRotator, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create log directory: %w", err)
    }

    fullPath := filepath.Join(logDir, logFileName)
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
        basePath:    logDir,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
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
    backupPath := filepath.Join(lr.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))
    currentPath := filepath.Join(lr.basePath, logFileName)

    if err := os.Rename(currentPath, backupPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    file, err := os.OpenFile(currentPath, os.O_CREATE|os.O_WRONLY, 0644)
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

    if len(matches) <= maxBackupFiles {
        return
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }

        timestampStr := parts[len(parts)-1]
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            continue
        }

        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{path: match, time: t})
    }

    if len(backupFiles) <= maxBackupFiles {
        return
    }

    for i := 0; i < len(backupFiles)-maxBackupFiles; i++ {
        os.Remove(backupFiles[i].path)
    }
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

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}