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
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    path := fmt.Sprintf("%s.log", basePath)
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentFile: file,
        filePath:    basePath,
        bytesWritten: info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.bytesWritten+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    rl.bytesWritten += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    oldPath := fmt.Sprintf("%s.log", rl.filePath)
    newPath := fmt.Sprintf("%s_%s.log", rl.filePath, timestamp)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    if err := rl.compressFile(newPath); err != nil {
        return err
    }

    file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.bytesWritten = 0
    rl.cleanupOldBackups()
    return nil
}

func (rl *RotatingLogger) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(path); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() error {
    pattern := fmt.Sprintf("%s_*.log.gz", rl.filePath)
    files, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(files) <= maxBackups {
        return nil
    }

    for i := 0; i < len(files)-maxBackups; i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: Application event occurred at %v\n", i, time.Now())
        if _, err := logger.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}