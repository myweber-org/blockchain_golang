
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const maxFileSize = 1024 * 1024 // 1MB
const maxBackups = 5

type LogRotator struct {
    filePath string
    currentSize int64
}

func NewLogRotator(path string) *LogRotator {
    return &LogRotator{
        filePath: path,
    }
}

func (lr *LogRotator) WriteLog(message string) error {
    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    info, err := file.Stat()
    if err != nil {
        return err
    }

    lr.currentSize = info.Size()
    newEntry := time.Now().Format("2006-01-02 15:04:05") + " - " + message + "\n"
    entrySize := int64(len(newEntry))

    if lr.currentSize+entrySize > maxFileSize {
        if err := lr.rotateLog(); err != nil {
            return err
        }
    }

    _, err = file.WriteString(newEntry)
    return err
}

func (lr *LogRotator) rotateLog() error {
    if err := lr.cleanOldBackups(); err != nil {
        return err
    }

    for i := maxBackups - 1; i > 0; i-- {
        oldName := lr.backupName(i)
        newName := lr.backupName(i + 1)
        
        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    if _, err := os.Stat(lr.filePath); err == nil {
        backupName := lr.backupName(1)
        return os.Rename(lr.filePath, backupName)
    }

    return nil
}

func (lr *LogRotator) backupName(index int) string {
    ext := filepath.Ext(lr.filePath)
    base := strings.TrimSuffix(lr.filePath, ext)
    return base + "." + strconv.Itoa(index) + ext
}

func (lr *LogRotator) cleanOldBackups() error {
    for i := maxBackups + 1; i <= maxBackups+10; i++ {
        backupFile := lr.backupName(i)
        if _, err := os.Stat(backupFile); err == nil {
            os.Remove(backupFile)
        }
    }
    return nil
}

func main() {
    rotator := NewLogRotator("application.log")
    
    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry number %d", i)
        if err := rotator.WriteLog(message); err != nil {
            fmt.Printf("Error writing log: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation demonstration completed")
}