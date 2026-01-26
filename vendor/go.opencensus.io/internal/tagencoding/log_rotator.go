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