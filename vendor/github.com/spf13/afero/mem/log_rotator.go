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

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentFile *os.File
    currentSize int64
    fileCounter int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    dir := filepath.Dir(l.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    l.currentFile = file
    l.currentSize = info.Size()
    l.fileCounter = 0
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()
    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }
    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", l.basePath, timestamp)
    if err := os.Rename(l.basePath, rotatedPath); err != nil {
        return err
    }
    if err := l.compressFile(rotatedPath); err != nil {
        return err
    }
    l.cleanOldFiles()
    return l.openCurrentFile()
}

func (l *RotatingLogger) compressFile(source string) error {
    dest := source + ".gz"
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    destFile, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer destFile.Close()
    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()
    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }
    os.Remove(source)
    return nil
}

func (l *RotatingLogger) cleanOldFiles() {
    dir := filepath.Dir(l.basePath)
    baseName := filepath.Base(l.basePath)
    files, err := os.ReadDir(dir)
    if err != nil {
        return
    }
    var compressedFiles []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            compressedFiles = append(compressedFiles, filepath.Join(dir, name))
        }
    }
    if len(compressedFiles) > 10 {
        filesToRemove := compressedFiles[:len(compressedFiles)-10]
        for _, file := range filesToRemove {
            os.Remove(file)
        }
    }
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("./logs/app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}
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

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
	currentSize int64
}

func NewRotatingLogger(filePath string, maxSizeMB int, backupCount int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	for i := rl.backupCount - 1; i >= 0; i-- {
		oldName := rl.getBackupName(i)
		newName := rl.getBackupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if i == rl.backupCount-1 {
				os.Remove(oldName)
			} else {
				os.Rename(oldName, newName)
			}
		}
	}

	firstBackup := rl.getBackupName(0)
	if err := os.Rename(rl.filePath, firstBackup); err != nil {
		return err
	}

	if err := rl.compressFile(firstBackup); err != nil {
		fmt.Printf("Compression failed: %v\n", err)
	}

	if err := rl.openCurrentFile(); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) getBackupName(index int) string {
	if index == 0 {
		return rl.filePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.filePath, index)
}

func (rl *RotatingLogger) compressFile(source string) error {
	dest := source + ".gz"

	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	os.Remove(source)
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs/app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}