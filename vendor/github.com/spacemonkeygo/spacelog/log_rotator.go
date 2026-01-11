package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupLogs = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    filePath    string
    currentSize int64
}

func NewRotatingLogger(baseDir string) (*RotatingLogger, error) {
    filePath := filepath.Join(baseDir, logFileName)
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
        currentFile: file,
        filePath:    filePath,
        currentSize: info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.currentSize+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.currentFile.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0

    go rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    dir := filepath.Dir(rl.filePath)
    pattern := filepath.Join(dir, logFileName+".*")

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackupLogs {
        filesToDelete := matches[:len(matches)-maxBackupLogs]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: This is a test log message for rotation testing", i)
        time.Sleep(10 * time.Millisecond)
    }
}