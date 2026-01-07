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
    file        *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    if maxFiles < 1 {
        maxFiles = 5
    }

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
        file:        file,
        basePath:    basePath,
        maxSize:     int64(maxSizeMB) * 1024 * 1024,
        currentSize: info.Size(),
        maxFiles:    maxFiles,
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

    for i := rl.maxFiles - 1; i > 0; i-- {
        oldPath := rl.getArchivePath(i)
        newPath := rl.getArchivePath(i + 1)

        if _, err := os.Stat(oldPath); err == nil {
            if i == rl.maxFiles-1 {
                os.Remove(newPath)
            } else {
                os.Rename(oldPath, newPath)
            }
        }
    }

    if err := os.Rename(rl.basePath, rl.getArchivePath(1)); err != nil {
        return err
    }

    if err := rl.compressArchive(1); err != nil {
        fmt.Printf("Compression failed: %v\n", err)
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    return nil
}

func (rl *RotatingLogger) getArchivePath(number int) string {
    if number == 0 {
        return rl.basePath
    }
    ext := filepath.Ext(rl.basePath)
    base := rl.basePath[:len(rl.basePath)-len(ext)]
    return fmt.Sprintf("%s.%d%s", base, number, ext)
}

func (rl *RotatingLogger) compressArchive(number int) error {
    sourcePath := rl.getArchivePath(number)
    destPath := sourcePath + ".gz"

    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(destPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}