
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type LogRotator struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    maxFiles    int
    currentSize int64
    currentFile *os.File
}

func NewLogRotator(basePath string, maxSizeMB int, maxFiles int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(lr.basePath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return err
    }

    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        err := lr.rotate()
        if err != nil {
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
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    err := os.Rename(lr.basePath, rotatedPath)
    if err != nil {
        return err
    }

    err = lr.compressFile(rotatedPath)
    if err != nil {
        return err
    }

    err = lr.cleanupOldFiles()
    if err != nil {
        return err
    }

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

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldFiles() error {
    dir := filepath.Dir(lr.basePath)
    baseName := filepath.Base(lr.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var compressedFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            compressedFiles = append(compressedFiles, filepath.Join(dir, name))
        }
    }

    if len(compressedFiles) <= lr.maxFiles {
        return nil
    }

    filesToDelete := len(compressedFiles) - lr.maxFiles
    for i := 0; i < filesToDelete; i++ {
        os.Remove(compressedFiles[i])
    }

    return nil
}

func (lr *LogRotator) extractTimestamp(filename string) (time.Time, error) {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return time.Time{}, fmt.Errorf("invalid filename format")
    }

    timestampStr := parts[len(parts)-2]
    return time.Parse("20060102_150405", timestampStr)
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
            time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}