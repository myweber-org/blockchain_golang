
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
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compressOld bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(lr.basePath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return fmt.Errorf("create directory failed: %w", err)
    }

    file, err := os.OpenFile(lr.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return fmt.Errorf("open log file failed: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return fmt.Errorf("stat log file failed: %w", err)
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        err := lr.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    err := os.Rename(lr.basePath, backupPath)
    if err != nil {
        return fmt.Errorf("rename log file failed: %w", err)
    }

    if lr.compressOld {
        compressedPath := backupPath + ".gz"
        err = compressFile(backupPath, compressedPath)
        if err != nil {
            fmt.Printf("compress failed: %v\n", err)
        } else {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }

    err = lr.cleanupOldBackups()
    if err != nil {
        fmt.Printf("cleanup failed: %v\n", err)
    }

    return lr.openCurrentFile()
}

func compressFile(src, dst string) error {
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

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    var backups []string
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") || isTimestampBackup(match, lr.basePath) {
            backups = append(backups, match)
        }
    }

    if len(backups) <= lr.maxBackups {
        return nil
    }

    toDelete := backups[:len(backups)-lr.maxBackups]
    for _, path := range toDelete {
        os.Remove(path)
    }

    return nil
}

func isTimestampBackup(path, basePath string) bool {
    suffix := strings.TrimPrefix(path, basePath+".")
    if len(suffix) != 14 {
        return false
    }
    _, err := strconv.ParseInt(suffix, 10, 64)
    return err == nil
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Write failed: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}