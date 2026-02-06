
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    filename   string
    current    *os.File
    size       int64
    backupList []string
}

func NewLogRotator(filename string) (*LogRotator, error) {
    lr := &LogRotator{filename: filename}
    if err := lr.openCurrent(); err != nil {
        return nil, err
    }
    lr.scanBackups()
    return lr, nil
}

func (lr *LogRotator) openCurrent() error {
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    lr.current = file
    lr.size = info.Size()
    return nil
}

func (lr *LogRotator) scanBackups() {
    pattern := lr.filename + ".*.gz"
    matches, _ := filepath.Glob(pattern)
    lr.backupList = matches
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.size+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := lr.current.Write(p)
    lr.size += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Unix()
    backupName := fmt.Sprintf("%s.%d", lr.filename, timestamp)
    if err := os.Rename(lr.filename, backupName); err != nil {
        return err
    }

    compressedName := backupName + ".gz"
    if err := compressFile(backupName, compressedName); err != nil {
        return err
    }
    os.Remove(backupName)

    lr.backupList = append(lr.backupList, compressedName)
    if len(lr.backupList) > maxBackups {
        oldest := lr.backupList[0]
        os.Remove(oldest)
        lr.backupList = lr.backupList[1:]
    }

    return lr.openCurrent()
}

func compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func (lr *LogRotator) parseBackupTimestamp(filename string) int64 {
    base := filepath.Base(filename)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return 0
    }
    ts, _ := strconv.ParseInt(parts[1], 10, 64)
    return ts
}

func (lr *LogRotator) Close() error {
    if lr.current != nil {
        return lr.current.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Println("Log rotation test completed")
}