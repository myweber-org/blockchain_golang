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
    filePath   string
    current    *os.File
    currentSize int64
}

func NewLogRotator(path string) (*LogRotator, error) {
    rotator := &LogRotator{filePath: path}
    if err := rotator.openCurrent(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.current.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.current != nil {
        lr.current.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)
    if err := os.Rename(lr.filePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := lr.cleanupOld(); err != nil {
        return err
    }

    return lr.openCurrent()
}

func (lr *LogRotator) compressFile(src string) error {
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

    gz := gzip.NewWriter(dstFile)
    defer gz.Close()

    if _, err := io.Copy(gz, srcFile); err != nil {
        return err
    }

    os.Remove(src)
    return nil
}

func (lr *LogRotator) cleanupOld() error {
    dir := filepath.Dir(lr.filePath)
    base := filepath.Base(lr.filePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, base+".") && strings.HasSuffix(name, ".gz") {
            backups = append(backups, name)
        }
    }

    if len(backups) <= maxBackups {
        return nil
    }

    for i := 0; i < len(backups)-maxBackups; i++ {
        os.Remove(filepath.Join(dir, backups[i]))
    }
    return nil
}

func (lr *LogRotator) openCurrent() error {
    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.current = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.current != nil {
        return lr.current.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log")
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}package main

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

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    maxFiles    int
    currentFile *os.File
    currentSize int64
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    l.currentFile = file
    l.currentSize = stat.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    for i := l.maxFiles - 1; i >= 0; i-- {
        oldPath := l.getArchivePath(i)
        if _, err := os.Stat(oldPath); os.IsNotExist(err) {
            continue
        }

        if i == l.maxFiles-1 {
            os.Remove(oldPath)
        } else {
            newPath := l.getArchivePath(i + 1)
            os.Rename(oldPath, newPath)
        }
    }

    if err := os.Rename(l.basePath, l.getArchivePath(0)); err != nil && !os.IsNotExist(err) {
        return err
    }

    return l.openCurrentFile()
}

func (l *RotatingLogger) getArchivePath(index int) string {
    if index == 0 {
        return l.basePath + ".1"
    }
    return l.basePath + "." + strconv.Itoa(index) + ".gz"
}

func (l *RotatingLogger) compressOldLogs() error {
    for i := 1; i < l.maxFiles; i++ {
        uncompressedPath := l.basePath + "." + strconv.Itoa(i)
        compressedPath := uncompressedPath + ".gz"

        if _, err := os.Stat(uncompressedPath); os.IsNotExist(err) {
            continue
        }

        if _, err := os.Stat(compressedPath); err == nil {
            continue
        }

        if err := compressFile(uncompressedPath, compressedPath); err != nil {
            return err
        }
        os.Remove(uncompressedPath)
    }
    return nil
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

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application event processed\n",
            time.Now().Format("2006-01-02 15:04:05"), i)
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    logger.compressOldLogs()
    fmt.Println("Log rotation completed")
}