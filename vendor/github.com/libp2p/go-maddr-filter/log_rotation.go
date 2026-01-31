package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    maxFiles    int
    currentFile *os.File
    currentSize int64
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    if maxFiles < 1 {
        return nil, fmt.Errorf("maxFiles must be at least 1")
    }
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) rotateIfNeeded() error {
    if rl.currentSize < rl.maxSize {
        return nil
    }

    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)
    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.openCurrentFile(); err != nil {
        return err
    }

    go rl.cleanupOldFiles()
    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        log.Printf("Failed to read directory for cleanup: %v", err)
        return
    }

    var rotatedFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && entry.Type().IsRegular() {
            rotatedFiles = append(rotatedFiles, filepath.Join(dir, name))
        }
    }

    if len(rotatedFiles) <= rl.maxFiles-1 {
        return
    }

    for i := 0; i < len(rotatedFiles)-rl.maxFiles+1; i++ {
        if err := os.Remove(rotatedFiles[i]); err != nil {
            log.Printf("Failed to remove old log file %s: %v", rotatedFiles[i], err)
        }
    }
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

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("./logs/app.log", 10, 5)
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    for i := 1; i <= 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(msg)); err != nil {
            log.Printf("Write error: %v", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
    fmt.Println("Log rotation demo completed. Check ./logs directory.")
}