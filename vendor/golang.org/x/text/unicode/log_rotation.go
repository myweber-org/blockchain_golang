package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
    logDir      = "logs"
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    basePath := filepath.Join(logDir, name)
    logger := &RotatingLogger{baseName: basePath}

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.baseName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
    backupPath := fmt.Sprintf("%s_%s.log", rl.baseName, timestamp)

    if err := os.Rename(rl.baseName+".log", backupPath); err != nil {
        return err
    }

    if err := rl.openCurrentFile(); err != nil {
        return err
    }

    rl.cleanupOldBackups()
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.baseName + "_*.log"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    filesToDelete := matches[:len(matches)-maxBackups]
    for _, file := range filesToDelete {
        os.Remove(file)
    }
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
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed")
}