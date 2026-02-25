
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type RotatingLogger struct {
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
}

func NewRotatingLogger(filePath string, maxSizeMB int) (*RotatingLogger, error) {
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
    
    return &RotatingLogger{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > rl.maxSize {
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
    rl.file.Close()
    
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    
    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    rl.file = file
    rl.currentSize = 0
    
    rl.cleanOldBackups()
    return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
    dir := filepath.Dir(rl.filePath)
    base := filepath.Base(rl.filePath)
    
    files, err := filepath.Glob(filepath.Join(dir, base+".*"))
    if err != nil {
        return
    }
    
    if len(files) > 5 {
        for i := 0; i < len(files)-5; i++ {
            os.Remove(files[i])
        }
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}