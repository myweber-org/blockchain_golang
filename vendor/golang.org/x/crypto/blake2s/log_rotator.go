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
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName    = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(logDir string) (*LogRotator, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    fullPath := filepath.Join(logDir, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
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
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    oldPath := filepath.Join(lr.basePath, logFileName)
    newPath := filepath.Join(lr.basePath, backupName)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
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

func (lr *LogRotator) parseBackupTimestamp(filename string) (time.Time, error) {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return time.Time{}, fmt.Errorf("invalid backup filename")
    }

    timestampStr := parts[len(parts)-1]
    return time.Parse("20060102_150405", timestampStr)
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
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}