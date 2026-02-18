
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentSize int64
    file        *os.File
    basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    fullPath := filepath.Join(path, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentSize: stat.Size(),
        file:        file,
        basePath:    path,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.currentSize+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.file.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.file.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s.gz", logFileName, timestamp)
    backupPath := filepath.Join(rl.basePath, backupName)

    oldLogPath := filepath.Join(rl.basePath, logFileName)
    oldFile, err := os.Open(oldLogPath)
    if err != nil {
        return err
    }
    defer oldFile.Close()

    backupFile, err := os.Create(backupPath)
    if err != nil {
        return err
    }
    defer backupFile.Close()

    gzWriter := gzip.NewWriter(backupFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }

    os.Remove(oldLogPath)

    newFile, err := os.Create(oldLogPath)
    if err != nil {
        return err
    }

    rl.file = newFile
    rl.currentSize = 0

    rl.cleanupOldBackups()
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := filepath.Join(rl.basePath, logFileName+".*.gz")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupFiles {
        return
    }

    for i := 0; i < len(matches)-maxBackupFiles; i++ {
        os.Remove(matches[i])
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(logger)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(10 * time.Millisecond)
    }
}