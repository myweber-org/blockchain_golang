package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logDir         = "./logs"
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    currentSize int64
    baseName   string
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    logger := &RotatingLogger{
        baseName: baseName,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    filePath := filepath.Join(logDir, rl.baseName+".log")
    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.file = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.file != nil {
        rl.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    oldPath := filepath.Join(logDir, rl.baseName+".log")
    newPath := filepath.Join(logDir, rl.baseName+"_"+timestamp+".log")

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    rl.cleanupOldFiles()
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
    pattern := filepath.Join(logDir, rl.baseName+"_*.log")
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

func (rl *RotatingLogger) parseBackupNumber(filename string) int {
    parts := strings.Split(filename, "_")
    if len(parts) < 2 {
        return -1
    }
    
    numStr := strings.TrimSuffix(parts[len(parts)-1], ".log")
    num, err := strconv.Atoi(numStr)
    if err != nil {
        return -1
    }
    return num
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(message))
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}