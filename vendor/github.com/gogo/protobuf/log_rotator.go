package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
)

type LogRotator struct {
    filename   string
    current   *os.File
    size      int64
}

func NewLogRotator(filename string) (*LogRotator, error) {
    lr := &LogRotator{filename: filename}
    if err := lr.openCurrent(); err != nil {
        return nil, err
    }
    return lr, nil
}

func (lr *LogRotator) openCurrent() error {
    file, err := os.OpenFile(lr.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    lr.current = file
    lr.size = info.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.size+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := lr.current.Write(p)
    lr.size += int64(n)
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.current.Close(); err != nil {
        return err
    }

    for i := maxBackups - 1; i >= 0; i-- {
        oldName := lr.backupName(i)
        newName := lr.backupName(i + 1)
        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    if err := os.Rename(lr.filename, lr.backupName(0)); err != nil {
        return err
    }

    return lr.openCurrent()
}

func (lr *LogRotator) backupName(index int) string {
    if index == 0 {
        return lr.filename + ".1"
    }
    return lr.filename + "." + strconv.Itoa(index+1)
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
        fmt.Fprintf(os.Stderr, "Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 10000; i++ {
        message := fmt.Sprintf("Log entry %d: This is a test log message.\n", i)
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
            break
        }
    }
    fmt.Println("Log rotation test completed")
}