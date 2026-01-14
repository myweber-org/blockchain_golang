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

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	
	dir := filepath.Dir(basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	logger := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
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
	archivePath := fmt.Sprintf("%s.%d.%s.gz", rl.filePath, rl.rotationCount, timestamp)
	rl.rotationCount++

	if err := compressFile(rl.filePath, archivePath); err != nil {
		return fmt.Errorf("failed to compress log file: %w", err)
	}

	if err := os.Remove(rl.filePath); err != nil {
		return fmt.Errorf("failed to remove old log file: %w", err)
	}

	return rl.openCurrentFile()
}

func compressFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}