package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) *LogRotator {
    return &LogRotator{
        basePath: basePath,
    }
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    file, err := os.OpenFile(filepath.Join(lr.basePath, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    n, err = file.Write(p)
    lr.currentSize += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    currentPath := filepath.Join(lr.basePath, logFileName)
    backupPath := filepath.Join(lr.basePath, logFileName+".1")

    if err := os.Rename(currentPath, backupPath); err != nil {
        return err
    }

    lr.currentSize = 0
    lr.cleanupOldBackups()
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    var backups []string
    pattern := filepath.Join(lr.basePath, logFileName+".*")

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    for _, match := range matches {
        if strings.HasSuffix(match, ".log") {
            backups = append(backups, match)
        }
    }

    sort.Slice(backups, func(i, j int) bool {
        numI := extractBackupNumber(backups[i])
        numJ := extractBackupNumber(backups[j])
        return numI > numJ
    })

    for i := maxBackupCount; i < len(backups); i++ {
        os.Remove(backups[i])
    }
}

func extractBackupNumber(path string) int {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    if len(parts) < 2 {
        return 0
    }
    num, _ := strconv.Atoi(parts[len(parts)-1])
    return num
}

func main() {
    rotator := NewLogRotator(".")
    writer := io.MultiWriter(os.Stdout, rotator)

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
            time.Now().Format("2006-01-02 15:04:05"), i)
        writer.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}