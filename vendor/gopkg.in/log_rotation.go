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

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    currentSize int64
    basePath   string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
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
        file:       file,
        currentSize: info.Size(),
        basePath:   path,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > maxFileSize {
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

    timestamp := time.Now().Unix()
    backupPath := fmt.Sprintf("%s.%d", rl.basePath, timestamp)
    if err := os.Rename(rl.basePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0

    go rl.compressBackup(backupPath)
    go rl.cleanOldBackups()

    return nil
}

func (rl *RotatingLogger) compressBackup(path string) {
    src, err := os.Open(path)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
    if err != nil {
        return
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return
    }

    os.Remove(path)
}

func (rl *RotatingLogger) cleanOldBackups() {
    pattern := filepath.Join(filepath.Dir(rl.basePath), filepath.Base(rl.basePath)+".*.gz")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    var backups []struct {
        path string
        time int64
    }

    for _, match := range matches {
        base := filepath.Base(match)
        timestampStr := base[len(filepath.Base(rl.basePath))+1 : len(base)-3]
        timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
        if err != nil {
            continue
        }
        backups = append(backups, struct {
            path string
            time int64
        }{match, timestamp})
    }

    for i := 0; i < len(backups)-maxBackups; i++ {
        oldestIdx := 0
        for j := 1; j < len(backups); j++ {
            if backups[j].time < backups[oldestIdx].time {
                oldestIdx = j
            }
        }
        os.Remove(backups[oldestIdx].path)
        backups = append(backups[:oldestIdx], backups[oldestIdx+1:]...)
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))))
        time.Sleep(10 * time.Millisecond)
    }
}