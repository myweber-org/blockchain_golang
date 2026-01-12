
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
    }
    
    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    fullPath := filepath.Join(rl.basePath, logFileName)
    
    file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    rl.currentFile = file
    rl.currentSize = info.Size()
    
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            log.Printf("Failed to rotate log: %v", err)
        }
    }
    
    n, err = rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupPath := filepath.Join(rl.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))
    
    currentPath := filepath.Join(rl.basePath, logFileName)
    if err := os.Rename(currentPath, backupPath); err != nil {
        return err
    }
    
    if err := rl.openCurrentFile(); err != nil {
        return err
    }
    
    go rl.cleanupOldBackups()
    
    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := filepath.Join(rl.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) <= maxBackupCount {
        return
    }
    
    var backups []struct {
        path string
        time time.Time
    }
    
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }
        
        timestampStr := parts[len(parts)-1]
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            continue
        }
        
        backups = append(backups, struct {
            path string
            time time.Time
        }{match, t})
    }
    
    for i := 0; i < len(backups)-maxBackupCount; i++ {
        oldestIdx := 0
        for j := 1; j < len(backups); j++ {
            if backups[j].time.Before(backups[oldestIdx].time) {
                oldestIdx = j
            }
        }
        
        os.Remove(backups[oldestIdx].path)
        backups = append(backups[:oldestIdx], backups[oldestIdx+1:]...)
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    log.SetOutput(io.MultiWriter(os.Stdout, logger))
    
    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: %s", i, strings.Repeat("x", 1024))
        time.Sleep(10 * time.Millisecond)
    }
}package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    sequence    int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: path,
        sequence: 0,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    if err := compressFile(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    rl.sequence++
    if rl.sequence > maxBackups {
        if err := rl.cleanOldBackups(); err != nil {
            return err
        }
    }

    return rl.openCurrentFile()
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

func (rl *RotatingLogger) cleanOldBackups() error {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) > maxBackups {
        toDelete := matches[:len(matches)-maxBackups]
        for _, file := range toDelete {
            if err := os.Remove(file); err != nil {
                return err
            }
        }
    }

    return nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
    logFileName = "app.log"
)

type LogRotator struct {
    currentSize int64
    file        *os.File
}

func NewLogRotator() (*LogRotator, error) {
    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        return nil, err
    }

    return &LogRotator{
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    if err := os.Rename(logFileName, backupName); err != nil {
        return err
    }

    file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0

    go lr.cleanupOldLogs()
    return nil
}

func (lr *LogRotator) cleanupOldLogs() {
    files, err := filepath.Glob(logFileName + ".*")
    if err != nil {
        return
    }

    if len(files) <= maxBackups {
        return
    }

    for i := 0; i < len(files)-maxBackups; i++ {
        os.Remove(files[i])
    }
}

func (lr *LogRotator) Close() error {
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator()
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}