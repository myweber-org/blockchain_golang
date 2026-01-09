
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu           sync.Mutex
    basePath     string
    currentFile  *os.File
    maxSize      int64
    currentSize  int64
    fileCount    int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    filename := rl.basePath
    if rl.fileCount > 0 {
        filename = fmt.Sprintf("%s.%d", rl.basePath, rl.fileCount)
    }

    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) rotate() error {
    rl.fileCount++
    return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func (rl *RotatingLogger) CleanupOldFiles(maxFiles int) error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    for i := 0; i <= rl.fileCount-maxFiles; i++ {
        filename := rl.basePath
        if i > 0 {
            filename = fmt.Sprintf("%s.%d", rl.basePath, i)
        }
        if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
            return err
        }
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 1024*1024)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    if err := logger.CleanupOldFiles(5); err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }
}