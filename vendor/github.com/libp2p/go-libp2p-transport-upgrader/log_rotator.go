
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

const (
    maxFileSize   = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type LogRotator struct {
    filePath    string
    currentFile *os.File
    currentSize int64
    mu          sync.Mutex
}

func NewLogRotator(path string) (*LogRotator, error) {
    lr := &LogRotator{
        filePath: path,
    }

    if err := lr.openCurrentFile(); err != nil {
        return nil, err
    }

    return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

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
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    backupFiles, err := lr.getBackupFiles()
    if err != nil {
        return err
    }

    if len(backupFiles) >= maxBackupFiles {
        oldest := backupFiles[len(backupFiles)-1]
        if err := os.Remove(oldest); err != nil {
            return err
        }
        backupFiles = backupFiles[:len(backupFiles)-1]
    }

    for i := len(backupFiles) - 1; i >= 0; i-- {
        oldPath := backupFiles[i]
        newPath := lr.generateBackupPath(i + 1)
        if err := os.Rename(oldPath, newPath); err != nil {
            return err
        }
    }

    compressedPath := lr.filePath + ".1.gz"
    if err := compressFile(lr.filePath, compressedPath); err != nil {
        return err
    }

    if err := os.Remove(lr.filePath); err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) getBackupFiles() ([]string, error) {
    dir := filepath.Dir(lr.filePath)
    base := filepath.Base(lr.filePath)

    files, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    var backups []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, base) && strings.HasSuffix(name, ".gz") {
            backups = append(backups, filepath.Join(dir, name))
        }
    }

    return backups, nil
}

func (lr *LogRotator) generateBackupPath(index int) string {
    return lr.filePath + "." + strconv.Itoa(index) + ".gz"
}

func compressFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
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
    rotator, err := NewLogRotator("app.log")
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
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}