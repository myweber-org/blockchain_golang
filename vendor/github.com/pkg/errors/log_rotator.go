
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
    mu          sync.Mutex
    file        *os.File
    currentSize int64
    maxSize     int64
    basePath    string
    maxBackups  int
}

func NewRotatingLog(basePath string, maxSize int64, maxBackups int) (*RotatingLog, error) {
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
        file:       file,
        currentSize: info.Size(),
        maxSize:     maxSize,
        basePath:    basePath,
        maxBackups:  maxBackups,
    }, nil
}

func (r *RotatingLog) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *RotatingLog) rotate() error {
    if err := r.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

    if err := os.Rename(r.basePath, backupPath); err != nil {
        return err
    }

    if err := r.compressBackup(backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(r.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    r.file = file
    r.currentSize = 0

    return r.cleanupOldBackups()
}

func (r *RotatingLog) compressBackup(srcPath string) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(srcPath + ".gz")
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    return os.Remove(srcPath)
}

func (r *RotatingLog) cleanupOldBackups() error {
    pattern := r.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= r.maxBackups {
        return nil
    }

    toDelete := matches[:len(matches)-r.maxBackups]
    for _, path := range toDelete {
        if err := os.Remove(path); err != nil {
            return err
        }
    }

    return nil
}

func (r *RotatingLog) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.file.Close()
}

func main() {
    log, err := NewRotatingLog("app.log", 1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(msg)); err != nil {
            panic(err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}