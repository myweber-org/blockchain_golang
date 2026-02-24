
package main

import (
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
    filePath string
    current  *os.File
    size     int64
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filePath: path}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
    file, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    rl.current = file
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := rl.current.Write(p)
    rl.size += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.current != nil {
        rl.current.Close()
    }

    for i := maxBackups - 1; i >= 0; i-- {
        oldPath := rl.backupPath(i)
        newPath := rl.backupPath(i + 1)
        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }

    if err := os.Rename(rl.filePath, rl.backupPath(0)); err != nil && !os.IsNotExist(err) {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLogger) backupPath(index int) string {
    if index == 0 {
        return rl.filePath + ".1"
    }
    return fmt.Sprintf("%s.%d", rl.filePath, index+1)
}

func (rl *RotatingLogger) Close() error {
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}