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

const (
    maxFileSize = 10 * 1024 * 1024
    maxBackups  = 5
)

type RotatingLogger struct {
    mu        sync.Mutex
    file      *os.File
    size      int64
    basePath  string
    currentID int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath:  path,
        currentID: 1,
    }

    if err := rl.openCurrent(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
    path := rl.currentFilePath()
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.file = file
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLogger) currentFilePath() string {
    return fmt.Sprintf("%s.%d.log", rl.basePath, rl.currentID)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    if err := rl.compressCurrent(); err != nil {
        return err
    }

    rl.currentID++
    if rl.currentID > maxBackups {
        rl.currentID = 1
    }

    return rl.openCurrent()
}

func (rl *RotatingLogger) compressCurrent() error {
    oldPath := rl.currentFilePath()
    newPath := oldPath + ".gz"

    src, err := os.Open(oldPath)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(newPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(oldPath); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application event processed\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}