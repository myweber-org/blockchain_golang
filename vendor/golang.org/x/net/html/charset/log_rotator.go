
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

type RotatingLogger struct {
    mu          sync.Mutex
    currentFile *os.File
    basePath    string
    maxSize     int64
    fileCount   int
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    stat, err := l.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if stat.Size()+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    return l.currentFile.Write(p)
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.gz", l.basePath, timestamp)

    if err := compressFile(l.basePath, rotatedPath); err != nil {
        return err
    }

    os.Remove(l.basePath)

    l.fileCount++
    if l.fileCount > l.maxFiles {
        l.cleanupOldFiles()
    }

    return l.openCurrentFile()
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

func (l *RotatingLogger) cleanupOldFiles() {
    pattern := l.basePath + ".*.gz"
    files, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(files) > l.maxFiles {
        filesToDelete := files[:len(files)-l.maxFiles]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    l.currentFile = file
    return nil
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().String())))
        time.Sleep(10 * time.Millisecond)
    }
}