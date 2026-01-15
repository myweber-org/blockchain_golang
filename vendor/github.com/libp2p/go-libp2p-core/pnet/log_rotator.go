
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
    mu         sync.Mutex
    basePath   string
    maxSize    int64
    current    *os.File
    currentSize int64
    maxFiles   int
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
    if maxFiles < 1 {
        maxFiles = 1
    }
    
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }
    
    if err := logger.openCurrent(); err != nil {
        return nil, err
    }
    
    return logger, nil
}

func (l *RotatingLogger) openCurrent() error {
    file, err := os.OpenFile(l.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    l.current = file
    l.currentSize = info.Size()
    return nil
}

func (l *RotatingLogger) rotate() error {
    l.current.Close()
    
    for i := l.maxFiles - 1; i > 0; i-- {
        oldPath := l.basePath + "." + strconv.Itoa(i)
        newPath := l.basePath + "." + strconv.Itoa(i+1)
        
        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }
    
    if _, err := os.Stat(l.basePath); err == nil {
        os.Rename(l.basePath, l.basePath+".1")
    }
    
    return l.openCurrent()
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := l.current.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    if l.current != nil {
        return l.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
        os.Exit(1)
    }
    defer logger.Close()
    
    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}