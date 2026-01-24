
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    currentFile   *os.File
    filePath      string
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewLogRotator(filePath string, maxSize int64) (*LogRotator, error) {
    rotator := &LogRotator{
        filePath: filePath,
        maxSize:  maxSize,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) rotateIfNeeded() error {
    if lr.currentSize < lr.maxSize {
        return nil
    }

    lr.currentFile.Close()
    lr.rotationCount++

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.%d", lr.filePath, timestamp, lr.rotationCount)
    
    if err := os.Rename(lr.filePath, archivedPath); err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if err := lr.rotateIfNeeded(); err != nil {
        return 0, err
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func (lr *LogRotator) CleanupOldLogs(maxAge time.Duration) error {
    dir := filepath.Dir(lr.filePath)
    baseName := filepath.Base(lr.filePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    cutoffTime := time.Now().Add(-maxAge)

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        info, err := entry.Info()
        if err != nil {
            continue
        }

        if info.ModTime().Before(cutoffTime) {
            matched, _ := filepath.Match(baseName+".*", entry.Name())
            if matched {
                oldPath := filepath.Join(dir, entry.Name())
                os.Remove(oldPath)
            }
        }
    }

    return nil
}