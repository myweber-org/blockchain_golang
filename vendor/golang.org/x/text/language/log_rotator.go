
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
    maxSize     int64
    basePath    string
    currentSize int64
    fileCount   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        maxSize:  maxSize,
        basePath: basePath,
    }
    
    if err := logger.rotateIfNeeded(); err != nil {
        return nil, err
    }
    
    return logger, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if err := rl.rotateIfNeeded(); err != nil {
        return 0, err
    }
    
    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
    if rl.currentFile == nil || rl.currentSize >= rl.maxSize {
        if rl.currentFile != nil {
            rl.currentFile.Close()
        }
        
        timestamp := time.Now().Format("20060102_150405")
        filename := fmt.Sprintf("%s_%d_%s.log", rl.basePath, rl.fileCount, timestamp)
        
        file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            return err
        }
        
        rl.currentFile = file
        rl.fileCount++
        
        if info, err := file.Stat(); err == nil {
            rl.currentSize = info.Size()
        } else {
            rl.currentSize = 0
        }
    }
    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app_log", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()
    
    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}