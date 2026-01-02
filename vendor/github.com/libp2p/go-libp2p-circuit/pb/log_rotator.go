package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentFile *os.File
    currentSize int64
    mu          sync.Mutex
}

func NewRotator(filePath string, maxSizeMB int, maxAgeHours int) (*Rotator, error) {
    absPath, err := filepath.Abs(filePath)
    if err != nil {
        return nil, err
    }

    dir := filepath.Dir(absPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    r := &Rotator{
        filePath: absPath,
        maxSize:  int64(maxSizeMB) * 1024 * 1024,
        maxAge:   time.Duration(maxAgeHours) * time.Hour,
    }

    if err := r.openCurrentFile(); err != nil {
        return nil, err
    }

    go r.cleanupOldFiles()
    return r, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.currentFile.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)
    if err := os.Rename(r.filePath, archivedPath); err != nil {
        return err
    }

    return r.openCurrentFile()
}

func (r *Rotator) openCurrentFile() error {
    file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) cleanupOldFiles() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        r.mu.Lock()
        cutoff := time.Now().Add(-r.maxAge)
        dir := filepath.Dir(r.filePath)
        base := filepath.Base(r.filePath)

        entries, err := os.ReadDir(dir)
        if err != nil {
            r.mu.Unlock()
            continue
        }

        for _, entry := range entries {
            if entry.IsDir() {
                continue
            }

            name := entry.Name()
            if len(name) <= len(base)+1 || name[:len(base)] != base {
                continue
            }

            info, err := entry.Info()
            if err != nil {
                continue
            }

            if info.ModTime().Before(cutoff) {
                oldPath := filepath.Join(dir, name)
                os.Remove(oldPath)
            }
        }
        r.mu.Unlock()
    }
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("./logs/app.log", 10, 24)
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