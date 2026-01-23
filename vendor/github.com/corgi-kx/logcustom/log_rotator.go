package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
)

type LogRotator struct {
    mu         sync.Mutex
    file       *os.File
    maxSize    int64
    basePath   string
    currentSize int64
    rotationCount int
}

func NewLogRotator(basePath string, maxSizeMB int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        file:       file,
        maxSize:    maxSize,
        basePath:   basePath,
        currentSize: info.Size(),
        rotationCount: 0,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    lr.rotationCount++
    archivedPath := fmt.Sprintf("%s.%d", lr.basePath, lr.rotationCount)
    if err := os.Rename(lr.basePath, archivedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0
    return nil
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("Log entry number %d\n", i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
    }

    fmt.Println("Log rotation test completed")
}