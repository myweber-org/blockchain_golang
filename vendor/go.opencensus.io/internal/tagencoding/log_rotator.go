package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type RotatingLogger struct {
    basePath   string
    maxSize    int64
    current    *os.File
    written    int64
}

func NewRotatingLogger(path string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return nil, err
    }
    
    return &RotatingLogger{
        basePath: path,
        maxSize:  maxSize,
        current:  f,
        written:  info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.written+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := rl.current.Write(p)
    if err == nil {
        rl.written += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.current.Close()
    
    timestamp := time.Now().Format("20060102_150405")
    archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)
    
    oldFile, err := os.Open(rl.basePath)
    if err != nil {
        return err
    }
    defer oldFile.Close()
    
    archiveFile, err := os.Create(archivePath)
    if err != nil {
        return err
    }
    defer archiveFile.Close()
    
    gzWriter := gzip.NewWriter(archiveFile)
    defer gzWriter.Close()
    
    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }
    
    os.Remove(rl.basePath)
    
    f, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    rl.current = f
    rl.written = 0
    return nil
}

func (rl *RotatingLogger) Close() error {
    return rl.current.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
    
    files, _ := filepath.Glob("app.log.*.gz")
    fmt.Printf("Created %d archived log files\n", len(files))
}
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
    currentFile   *os.File
    currentSize   int64
    basePath      string
    currentSuffix int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    if err := rotator.openCurrentFile(); err != nil {
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

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    lr.cleanupOldBackups()

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    destFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    backupFiles := make([]backupFile, 0, len(matches))
    for _, match := range matches {
        if ts, err := extractTimestamp(match); err == nil {
            backupFiles = append(backupFiles, backupFile{path: match, timestamp: ts})
        }
    }

    sortBackupsByTime(backupFiles)

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        os.Remove(backupFiles[i].path)
    }
}

func (lr *LogRotator) openCurrentFile() error {
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

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

type backupFile struct {
    path      string
    timestamp time.Time
}

func extractTimestamp(path string) (time.Time, error) {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return time.Time{}, fmt.Errorf("invalid backup file name")
    }

    timestampStr := parts[1]
    return time.Parse("20060102_150405", timestampStr)
}

func sortBackupsByTime(files []backupFile) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if files[i].timestamp.After(files[j].timestamp) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
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