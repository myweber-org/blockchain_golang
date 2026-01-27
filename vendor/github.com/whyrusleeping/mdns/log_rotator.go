package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logFileName = "app.log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator(path string) (*LogRotator, error) {
	fullPath := filepath.Join(path, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		basePath:    path,
	}, nil
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
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	oldPath := filepath.Join(lr.basePath, logFileName)
	newPath := filepath.Join(lr.basePath, backupName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0

	go lr.cleanupOldBackups()

	return nil
}

func (lr *LogRotator) cleanupOldBackups() {
	pattern := filepath.Join(lr.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackups {
		return
	}

	for i := 0; i < len(matches)-maxBackups; i++ {
		os.Remove(matches[i])
	}
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator(".")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		rotator.Write([]byte(message))
		time.Sleep(100 * time.Millisecond)
	}
}package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type LogRotator struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    rotation    int
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }
    
    return &LogRotator{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
        rotation:    0,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    
    if lr.currentSize+int64(len(p)) > lr.maxSize {
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
    if lr.file != nil {
        lr.file.Close()
    }
    
    timestamp := time.Now().Format("20060102_150405")
    archivedName := fmt.Sprintf("%s.%s.%d.gz", lr.filePath, timestamp, lr.rotation)
    
    source, err := os.Open(lr.filePath)
    if err != nil {
        return err
    }
    defer source.Close()
    
    dest, err := os.Create(archivedName)
    if err != nil {
        return err
    }
    defer dest.Close()
    
    gzWriter := gzip.NewWriter(dest)
    defer gzWriter.Close()
    
    if _, err := io.Copy(gzWriter, source); err != nil {
        return err
    }
    
    if err := os.Remove(lr.filePath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    lr.file = file
    lr.currentSize = 0
    lr.rotation++
    
    lr.cleanOldArchives()
    return nil
}

func (lr *LogRotator) cleanOldArchives() {
    pattern := filepath.Base(lr.filePath) + ".*.gz"
    matches, err := filepath.Glob(filepath.Join(filepath.Dir(lr.filePath), pattern))
    if err != nil {
        return
    }
    
    if len(matches) > 10 {
        for i := 0; i < len(matches)-10; i++ {
            os.Remove(matches[i])
        }
    }
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    
    if lr.file != nil {
        return lr.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log", 10)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation completed")
}