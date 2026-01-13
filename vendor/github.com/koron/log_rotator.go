
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
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logDir        = "./logs"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    baseName    string
    fileIndex   int
}

func NewLogRotator(baseName string) (*LogRotator, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    rotator := &LogRotator{
        baseName: baseName,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    path := filepath.Join(logDir, lr.baseName+".log")
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedName := fmt.Sprintf("%s_%s.log", lr.baseName, timestamp)
    oldPath := filepath.Join(logDir, lr.baseName+".log")
    newPath := filepath.Join(logDir, rotatedName)

    if err := os.Rename(oldPath, newPath); err != nil {
        return err
    }

    if err := lr.compressFile(newPath); err != nil {
        return err
    }

    lr.cleanupOldFiles()

    lr.currentSize = 0
    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := filepath.Join(logDir, lr.baseName+"_*.log.gz")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupFiles {
        return
    }

    var files []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(filepath.Base(match), "_")
        if len(parts) < 3 {
            continue
        }

        timestampStr := parts[1] + "_" + strings.TrimSuffix(parts[2], ".log.gz")
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            continue
        }

        files = append(files, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(files)-maxBackupFiles; i++ {
        oldestIndex := i
        for j := i + 1; j < len(files); j++ {
            if files[j].time.Before(files[oldestIndex].time) {
                oldestIndex = j
            }
        }
        os.Remove(files[oldestIndex].path)
        files[oldestIndex] = files[i]
    }
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }

        if i%100 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Println("Log rotation test completed")
}