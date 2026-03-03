package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu         sync.Mutex
    basePath   string
    maxSize    int64
    file       *os.File
    currentSize int64
    maxBackups int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxBackups int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        basePath:   basePath,
        maxSize:    maxSize,
        file:       file,
        currentSize: info.Size(),
        maxBackups: maxBackups,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0

    return rl.cleanupOldBackups()
}

func (rl *RotatingLogger) compressFile(src string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(src + ".gz")
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    return os.Remove(src)
}

func (rl *RotatingLogger) cleanupOldBackups() error {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= rl.maxBackups {
        return nil
    }

    for i := 0; i < len(matches)-rl.maxBackups; i++ {
        if err := os.Remove(matches[i]); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLogger) parseBackupNumber(filename string) int {
    base := filepath.Base(filename)
    if len(base) <= len(rl.basePath) {
        return 0
    }

    suffix := base[len(rl.basePath)+1:]
    numStr := ""
    for i := 0; i < len(suffix); i++ {
        if suffix[i] >= '0' && suffix[i] <= '9' {
            numStr += string(suffix[i])
        } else {
            break
        }
    }

    if num, err := strconv.Atoi(numStr); err == nil {
        return num
    }
    return 0
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}