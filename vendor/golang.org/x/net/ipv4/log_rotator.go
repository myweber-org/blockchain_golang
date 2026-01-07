
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
    }
    
    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    fullPath := filepath.Join(rl.basePath, logFileName)
    
    file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    rl.currentFile = file
    rl.currentSize = info.Size()
    
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            log.Printf("Failed to rotate log: %v", err)
        }
    }
    
    n, err = rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupPath := filepath.Join(rl.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))
    
    currentPath := filepath.Join(rl.basePath, logFileName)
    if err := os.Rename(currentPath, backupPath); err != nil {
        return err
    }
    
    if err := rl.openCurrentFile(); err != nil {
        return err
    }
    
    go rl.cleanupOldBackups()
    
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := filepath.Join(rl.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) <= maxBackupCount {
        return
    }
    
    var backups []struct {
        path string
        time time.Time
    }
    
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }
        
        timestampStr := parts[len(parts)-1]
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            continue
        }
        
        backups = append(backups, struct {
            path string
            time time.Time
        }{match, t})
    }
    
    for i := 0; i < len(backups)-maxBackupCount; i++ {
        oldestIdx := 0
        for j := 1; j < len(backups); j++ {
            if backups[j].time.Before(backups[oldestIdx].time) {
                oldestIdx = j
            }
        }
        
        os.Remove(backups[oldestIdx].path)
        backups = append(backups[:oldestIdx], backups[oldestIdx+1:]...)
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    log.SetOutput(io.MultiWriter(os.Stdout, logger))
    
    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: %s", i, strings.Repeat("x", 1024))
        time.Sleep(10 * time.Millisecond)
    }
}