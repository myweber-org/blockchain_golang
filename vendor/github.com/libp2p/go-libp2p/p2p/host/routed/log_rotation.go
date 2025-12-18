package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    currentSize int64
    basePath   string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("create directory: %w", err)
    }

    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("stat log file: %w", err)
    }

    return &RotatingLogger{
        file:       file,
        currentSize: info.Size(),
        basePath:   path,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, fmt.Errorf("rotate log file: %w", err)
        }
    }

    n, err := rl.file.Write(p)
    if err != nil {
        return n, fmt.Errorf("write to log file: %w", err)
    }

    rl.currentSize += int64(n)
    return n, nil
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return fmt.Errorf("close current log file: %w", err)
    }

    timestamp := time.Now().Format("20060102-150405")
    backupPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, backupPath); err != nil {
        return fmt.Errorf("rename log file: %w", err)
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("create new log file: %w", err)
    }

    rl.file = file
    rl.currentSize = 0

    go rl.cleanupOldBackups()
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, path := range toDelete {
            os.Remove(path)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("logs/application.log")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}