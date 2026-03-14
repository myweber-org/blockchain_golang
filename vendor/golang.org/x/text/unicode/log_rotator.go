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

	logger := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
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
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == rl.backupCount-1 {
				os.Remove(oldPath)
			} else {
				if err := rl.compressFile(oldPath, newPath); err != nil {
					return err
				}
				os.Remove(oldPath)
			}
		}
	}

	firstBackup := rl.getBackupPath(0)
	if err := rl.compressFile(rl.filePath, firstBackup); err != nil {
		return err
	}

	os.Remove(rl.filePath)
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src, dst string) error {
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

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.filePath + ".1.gz"
	}
	return fmt.Sprintf("%s.%d.gz", rl.filePath, index+1)
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
		message := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
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
    basePath    string
    maxSize     int64
    fileSize    int64
    fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath:  basePath,
        maxSize:   maxSize,
        fileIndex: 0,
    }
    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", l.basePath, l.fileIndex)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    l.currentFile = file
    l.fileSize = stat.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.fileSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := l.currentFile.Write(p)
    if err == nil {
        l.fileSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
        go l.compressPreviousFile()
    }
    l.fileIndex++
    return l.openCurrentFile()
}

func (l *RotatingLogger) compressPreviousFile() {
    prevIndex := l.fileIndex - 1
    if prevIndex < 0 {
        return
    }

    srcName := fmt.Sprintf("%s.%d.log", l.basePath, prevIndex)
    dstName := fmt.Sprintf("%s.%d.log.gz", l.basePath, prevIndex)

    srcFile, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return
    }

    os.Remove(srcName)
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
    logger, err := NewRotatingLogger("app", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application event processed\n",
            time.Now().Format("2006-01-02 15:04:05"), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}