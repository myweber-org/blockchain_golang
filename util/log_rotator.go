
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingWriter struct {
    filename   string
    current    *os.File
    size       int64
    mu         sync.Mutex
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
    w := &RotatingWriter{filename: filename}
    if err := w.openFile(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.size+int64(len(p)) >= maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.current.Write(p)
    w.size += int64(n)
    return n, err
}

func (w *RotatingWriter) openFile() error {
    file, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    w.current = file
    w.size = info.Size()
    return nil
}

func (w *RotatingWriter) rotate() error {
    if w.current != nil {
        w.current.Close()
    }

    timestamp := time.Now().Format("20060102-150405")
    backupName := fmt.Sprintf("%s.%s.gz", w.filename, timestamp)
    
    if err := compressFile(w.filename, backupName); err != nil {
        return err
    }

    if err := cleanupOldBackups(w.filename); err != nil {
        return err
    }

    os.Remove(w.filename)
    return w.openFile()
}

func compressFile(source, target string) error {
    src, err := os.Open(source)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    return err
}

func cleanupOldBackups(baseName string) error {
    pattern := fmt.Sprintf("%s.*.gz", filepath.Base(baseName))
    matches, err := filepath.Glob(filepath.Join(filepath.Dir(baseName), pattern))
    if err != nil {
        return err
    }

    if len(matches) > maxBackups {
        filesToDelete := matches[:len(matches)-maxBackups]
        for _, f := range filesToDelete {
            os.Remove(f)
        }
    }
    return nil
}

func (w *RotatingWriter) Close() error {
    w.mu.Lock()
    defer w.mu.Unlock()
    if w.current != nil {
        return w.current.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log")
    if err != nil {
        panic(err)
    }
    defer writer.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
}