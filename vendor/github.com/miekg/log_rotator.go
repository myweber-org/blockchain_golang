
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
)

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCounter int
}

func NewRotatingLog(basePath string, maxSize int64) (*RotatingLog, error) {
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

    rl.fileCounter++
    archivePath := rl.basePath + "." + strconv.Itoa(rl.fileCounter) + ".gz"

    if err := compressFile(rl.basePath, archivePath); err != nil {
        return err
    }

    if err := os.Truncate(rl.basePath, 0); err != nil {
        return err
    }

    return rl.openCurrentFile()
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
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
    logDir := "./logs"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        fmt.Printf("Failed to create log directory: %v\n", err)
        return
    }

    logPath := filepath.Join(logDir, "application.log")
    rotator, err := NewRotatingLog(logPath, 1024*1024)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: This is a sample log message for testing rotation\n", i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
    }

    fmt.Println("Log rotation test completed. Check ./logs directory for rotated files.")
}