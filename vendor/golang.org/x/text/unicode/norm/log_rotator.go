
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type RotatingLog struct {
    currentFile   *os.File
    basePath      string
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLog{
        currentFile:   file,
        basePath:      basePath,
        maxSize:       maxSize,
        currentSize:   info.Size(),
        rotationCount: 0,
    }, nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
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
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    if err := compressFile(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    rl.rotationCount++

    return nil
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLog) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logDir := "./logs"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        panic(err)
    }

    logPath := filepath.Join(logDir, "application.log")
    rotator, err := NewRotatingLog(logPath, 10)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}