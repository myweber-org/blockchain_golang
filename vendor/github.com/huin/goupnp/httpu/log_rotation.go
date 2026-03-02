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
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
    logDir      = "./logs"
)

type RotatingLogger struct {
    currentFile *os.File
    baseName    string
    fileIndex   int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        baseName: baseName,
    }

    if err := rl.openNextFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openNextFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    fileName := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.fileIndex))
    file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.fileIndex++

    return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    fileInfo, err := rl.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if fileInfo.Size()+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.openNextFile(); err != nil {
        return err
    }

    return rl.cleanupOldFiles()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
    files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
    if err != nil {
        return err
    }

    if len(files) <= maxBackups {
        return nil
    }

    var fileIndices []int
    for _, file := range files {
        base := filepath.Base(file)
        idxStr := strings.TrimSuffix(strings.TrimPrefix(base, rl.baseName+"_"), ".log")
        idx, err := strconv.Atoi(idxStr)
        if err != nil {
            continue
        }
        fileIndices = append(fileIndices, idx)
    }

    for i := 0; i < len(fileIndices)-maxBackups; i++ {
        oldFile := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, fileIndices[i]))
        if err := os.Remove(oldFile); err != nil {
            log.Printf("Failed to remove old log file %s: %v", oldFile, err)
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
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 100; i++ {
        log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}