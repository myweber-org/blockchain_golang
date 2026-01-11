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

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    fileIndex   int
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
    w := &RotatingWriter{
        basePath:  path,
        fileIndex: 0,
    }

    if err := w.openCurrentFile(); err != nil {
        return nil, err
    }

    return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.currentFile.Write(p)
    w.currentSize += int64(n)
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if err := w.currentFile.Close(); err != nil {
        return err
    }

    if err := w.compressCurrentFile(); err != nil {
        return err
    }

    w.fileIndex++
    if w.fileIndex > maxBackups {
        w.cleanOldBackups()
    }

    return w.openCurrentFile()
}

func (w *RotatingWriter) compressCurrentFile() error {
    oldPath := w.currentFile.Name()
    compressedPath := oldPath + ".gz"

    src, err := os.Open(oldPath)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    if err != nil {
        return err
    }

    return os.Remove(oldPath)
}

func (w *RotatingWriter) openCurrentFile() error {
    timestamp := time.Now().Format("2006-01-02")
    filename := fmt.Sprintf("%s-%s.log", w.basePath, timestamp)
    
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    w.currentFile = file
    w.currentSize = stat.Size()
    return nil
}

func (w *RotatingWriter) cleanOldBackups() {
    pattern := w.basePath + "-*.log.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(filepath.Base(match), "-")
        if len(parts) < 2 {
            continue
        }

        dateStr := strings.TrimSuffix(parts[1], ".log.gz")
        t, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            continue
        }

        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        os.Remove(backupFiles[i].path)
    }
}

func (w *RotatingWriter) Close() error {
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("application")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        return
    }
    defer writer.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n", 
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}