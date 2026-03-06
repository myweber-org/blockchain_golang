
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type Rotator struct {
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentSize int64
    file        *os.File
}

func NewRotator(filePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    r := &Rotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxAge:   maxAge,
    }

    if err := r.openFile(); err != nil {
        return nil, err
    }

    go r.ageMonitor()
    return r, nil
}

func (r *Rotator) openFile() error {
    if err := os.MkdirAll(filepath.Dir(r.filePath), 0755); err != nil {
        return err
    }

    info, err := os.Stat(r.filePath)
    if os.IsNotExist(err) {
        file, err := os.Create(r.filePath)
        if err != nil {
            return err
        }
        r.file = file
        r.currentSize = 0
        return nil
    }
    if err != nil {
        return err
    }

    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    r.file = file
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err != nil {
        return n, err
    }

    r.currentSize += int64(n)
    return n, nil
}

func (r *Rotator) rotate() error {
    if r.file != nil {
        r.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }

    return r.openFile()
}

func (r *Rotator) ageMonitor() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        info, err := os.Stat(r.filePath)
        if err != nil {
            continue
        }

        if time.Since(info.ModTime()) > r.maxAge {
            r.rotate()
        }
    }
}

func (r *Rotator) Close() error {
    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator(
        "/var/log/myapp/app.log",
        10*1024*1024,
        24*time.Hour,
    )
    if err != nil {
        fmt.Printf("Failed to create rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Printf("Write failed: %v\n", err)
            break
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation example completed")
}