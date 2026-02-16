
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLog struct {
    mu            sync.Mutex
    currentFile   *os.File
    basePath      string
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  int64(maxSizeMB) * 1024 * 1024,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) openCurrentFile() error {
    f, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }

    rl.currentFile = f
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
    }

    timestamp := time.Now().Format("20060102_150405")
    archivePath := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.rotationCount, timestamp)

    if err := compressFile(rl.basePath, archivePath); err != nil {
        return fmt.Errorf("compression failed: %w", err)
    }

    if err := os.Remove(rl.basePath); err != nil {
        return fmt.Errorf("failed to remove original file: %w", err)
    }

    rl.rotationCount++
    return rl.openCurrentFile()
}

func compressFile(source, target string) error {
    src, err := os.Open(source)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
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
    logFile := "application.log"
    rotator, err := NewRotatingLog(logFile, 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Processing request from client\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed. Check", logFile, "and compressed archives.")
}