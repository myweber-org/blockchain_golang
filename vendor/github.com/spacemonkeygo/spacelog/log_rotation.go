package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logDirectory = "./logs"
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    baseName    string
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDirectory, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        baseName: baseName,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    fileName := filepath.Join(logDirectory, rl.baseName+".log")
    file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
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

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    oldPath := filepath.Join(logDirectory, rl.baseName+".log")
    newPath := filepath.Join(logDirectory, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    if err := rl.cleanupOldFiles(); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
    pattern := filepath.Join(logDirectory, rl.baseName+"_*.log")
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) <= maxBackups {
        return nil
    }

    for i := 0; i < len(files)-maxBackups; i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}