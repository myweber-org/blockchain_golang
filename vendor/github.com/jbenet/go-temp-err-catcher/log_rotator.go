package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type LogRotator struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }
    
    return &LogRotator{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
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
    if lr.file != nil {
        lr.file.Close()
    }
    
    timestamp := time.Now().Format("20060102_150405")
    dir := filepath.Dir(lr.filePath)
    base := filepath.Base(lr.filePath)
    ext := filepath.Ext(base)
    name := base[:len(base)-len(ext)]
    
    archivePath := filepath.Join(dir, fmt.Sprintf("%s_%s%s", name, timestamp, ext))
    
    if err := os.Rename(lr.filePath, archivePath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
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
    
    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}