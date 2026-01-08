package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    currentFile *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath:  basePath,
        maxSize:   maxSize,
        fileIndex: 0,
    }
    
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", l.basePath, l.fileIndex)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    l.currentFile = file
    l.currentSize = info.Size()
    return nil
}

func (l *RotatingLogger) rotateIfNeeded(data []byte) error {
    if l.currentSize+int64(len(data)) > l.maxSize {
        l.currentFile.Close()
        l.fileIndex++
        return l.openCurrentFile()
    }
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    if err := l.rotateIfNeeded(p); err != nil {
        return 0, err
    }
    
    n, err := l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("application", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()
    
    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(message))
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}