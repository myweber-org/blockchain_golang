
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCounter int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  maxSize,
    }
    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLog) openCurrentFile() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLog) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
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

func (rl *RotatingLog) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
        rl.currentFile = nil
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.%d", rl.basePath, timestamp, rl.fileCounter)
    rl.fileCounter++

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    go rl.compressOldLog(rotatedPath)

    return rl.openCurrentFile()
}

func (rl *RotatingLog) compressOldLog(path string) {
    originalFile, err := os.Open(path)
    if err != nil {
        return
    }
    defer originalFile.Close()

    compressedPath := path + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, originalFile); err != nil {
        return
    }

    os.Remove(path)
}

func (rl *RotatingLog) findOldLogs() []string {
    pattern := rl.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return nil
    }

    var logs []string
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") {
            logs = append(logs, match)
        }
    }
    return logs
}

func (rl *RotatingLog) cleanupOldLogs(maxFiles int) {
    logs := rl.findOldLogs()
    if len(logs) <= maxFiles {
        return
    }

    logsToDelete := logs[:len(logs)-maxFiles]
    for _, log := range logsToDelete {
        os.Remove(log)
    }
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    log.cleanupOldLogs(5)
    fmt.Println("Log rotation completed")
}