package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type Rotator struct {
    FilePath    string
    MaxSize     int64
    MaxAge      time.Duration
    currentFile *os.File
    currentSize int64
}

func NewRotator(filePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    r := &Rotator{
        FilePath: filePath,
        MaxSize:  maxSize,
        MaxAge:   maxAge,
    }
    if err := r.openCurrent(); err != nil {
        return nil, err
    }
    go r.timeBasedRotation()
    return r, nil
}

func (r *Rotator) openCurrent() error {
    if err := os.MkdirAll(filepath.Dir(r.FilePath), 0755); err != nil {
        return err
    }
    f, err := os.OpenFile(r.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    r.currentFile = f
    r.currentSize = stat.Size()
    return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    if r.currentSize >= r.MaxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := r.currentFile.Write(p)
    if err != nil {
        return n, err
    }
    r.currentSize += int64(n)
    return n, nil
}

func (r *Rotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.FilePath, timestamp)
    if err := os.Rename(r.FilePath, backupPath); err != nil {
        return err
    }
    return r.openCurrent()
}

func (r *Rotator) timeBasedRotation() {
    ticker := time.NewTicker(r.MaxAge)
    defer ticker.Stop()
    for range ticker.C {
        if r.currentSize > 0 {
            r.rotate()
        }
    }
}

func (r *Rotator) Close() error {
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("/var/log/app/app.log", 10*1024*1024, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}