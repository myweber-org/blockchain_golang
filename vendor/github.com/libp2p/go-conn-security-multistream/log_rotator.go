
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

type LogRotator struct {
    mu          sync.Mutex
    currentFile *os.File
    filePath    string
    maxSize     int64
    maxFiles    int
}

func NewLogRotator(filePath string, maxSizeMB int, maxFiles int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    rotator := &LogRotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    stat, err := lr.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if stat.Size()+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    return lr.currentFile.Write(p)
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", lr.filePath, timestamp)

    if err := compressFile(lr.filePath, archivedPath); err != nil {
        return err
    }

    if err := os.Remove(lr.filePath); err != nil {
        return err
    }

    if err := lr.cleanupOldFiles(); err != nil {
        return err
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

func (lr *LogRotator) cleanupOldFiles() error {
    pattern := fmt.Sprintf("%s.*.gz", lr.filePath)
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) <= lr.maxFiles {
        return nil
    }

    filesToDelete := files[:len(files)-lr.maxFiles]
    for _, file := range filesToDelete {
        if err := os.Remove(file); err != nil {
            return err
        }
    }
    return nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    lr.currentFile = file
    return nil
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}