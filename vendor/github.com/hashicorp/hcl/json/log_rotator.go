
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err != nil {
        return n, err
    }

    lr.currentSize += int64(n)
    return n, nil
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    backupMap := make(map[int64]string)
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }

        timestampStr := parts[len(parts)-2]
        timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
        if err != nil {
            continue
        }

        backupMap[timestamp] = match
    }

    timestamps := make([]int64, 0, len(backupMap))
    for ts := range backupMap {
        timestamps = append(timestamps, ts)
    }

    for i := 0; i < len(timestamps)-maxBackups; i++ {
        oldest := timestamps[i]
        os.Remove(backupMap[oldest])
    }

    return nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}