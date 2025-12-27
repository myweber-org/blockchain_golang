
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

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
}

func NewLogRotator(basePath string, maxSize int64, maxBackups int) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath:   basePath,
        maxSize:    maxSize,
        maxBackups: maxBackups,
    }

    if err := rotator.openCurrent(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
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

func (r *LogRotator) openCurrent() error {
    dir := filepath.Dir(r.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    info, err := os.Stat(r.basePath)
    if os.IsNotExist(err) {
        file, err := os.Create(r.basePath)
        if err != nil {
            return err
        }
        r.currentFile = file
        r.currentSize = 0
        return nil
    }
    if err != nil {
        return err
    }

    file, err := os.OpenFile(r.basePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

    if err := os.Rename(r.basePath, rotatedPath); err != nil {
        return err
    }

    if err := r.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := r.cleanupOldBackups(); err != nil {
        return err
    }

    return r.openCurrent()
}

func (r *LogRotator) compressFile(src string) error {
    dst := src + ".gz"
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    os.Remove(src)
    return nil
}

func (r *LogRotator) cleanupOldBackups() error {
    pattern := r.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= r.maxBackups {
        return nil
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(filepath.Base(match), ".")
        if len(parts) < 3 {
            continue
        }
        timestamp := parts[len(parts)-2]
        t, err := time.Parse("20060102150405", timestamp)
        if err != nil {
            continue
        }
        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(backupFiles)-r.maxBackups; i++ {
        os.Remove(backupFiles[i].path)
    }

    return nil
}

func (r *LogRotator) extractTimestamp(filename string) (time.Time, error) {
    base := filepath.Base(filename)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return time.Time{}, fmt.Errorf("invalid filename format")
    }
    timestamp := parts[len(parts)-2]
    return time.Parse("20060102150405", timestamp)
}

func (r *LogRotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10*1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample data here\n",
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation completed")
}