
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentSize int64
    file        *os.File
}

func NewLogRotator() (*LogRotator, error) {
    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxLogSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    lr.currentSize += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    lr.file.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s.gz", logFileName, timestamp)

    if err := compressFile(logFileName, backupName); err != nil {
        return err
    }

    if err := cleanupOldBackups(); err != nil {
        return err
    }

    file, err := os.Create(logFileName)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0
    return nil
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

func cleanupOldBackups() error {
    pattern := fmt.Sprintf("%s.*.gz", logFileName)
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) <= maxBackupFiles {
        return nil
    }

    for i := 0; i < len(files)-maxBackupFiles; i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }

    return nil
}

func (lr *LogRotator) Close() error {
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator()
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}