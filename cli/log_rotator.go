package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
)

type RotatingLogger struct {
    currentFile *os.File
    basePath    string
    fileSize    int64
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    rl := &RotatingLogger{basePath: basePath}
    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    filename := rl.basePath + logExtension
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.fileSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.fileSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.fileSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s_%s%s", rl.basePath, timestamp, logExtension)

    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    if err := os.Rename(rl.basePath+logExtension, rotatedPath); err != nil {
        return err
    }

    if err := rl.openCurrentFile(); err != nil {
        return err
    }

    rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    files, err := os.ReadDir(dir)
    if err != nil {
        log.Printf("Failed to read directory: %v", err)
        return
    }

    var logFiles []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, baseName) && strings.HasSuffix(name, logExtension) && name != filepath.Base(rl.basePath+logExtension) {
            logFiles = append(logFiles, filepath.Join(dir, name))
        }
    }

    if len(logFiles) <= maxBackups {
        return
    }

    for i := 0; i < len(logFiles)-maxBackups; i++ {
        if err := os.Remove(logFiles[i]); err != nil {
            log.Printf("Failed to remove old log file %s: %v", logFiles[i], err)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app_log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    multiWriter := io.MultiWriter(os.Stdout, logger)
    log.SetOutput(multiWriter)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running smoothly", i)
        time.Sleep(10 * time.Millisecond)
    }
}