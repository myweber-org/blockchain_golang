package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type LogRotator struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
}

func NewLogRotator(filePath string, maxSize int64) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	backupPath := fmt.Sprintf("%s.1", lr.filePath)
	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) Write(data []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(data)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(data)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
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
	rotator, err := NewLogRotator("app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i)
		_, err := rotator.Write([]byte(message))
		if err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     maxSize,
		backupCount: backupCount,
		compressOld: compressOld,
	}

	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile == nil {
		if err := rl.openCurrentFile(); err != nil {
			return 0, err
		}
	}

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	rl.currentSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) rotateIfNeeded() error {
	info, err := os.Stat(rl.basePath)
	if os.IsNotExist(err) {
		return rl.openCurrentFile()
	}
	if err != nil {
		return err
	}

	if info.Size() >= rl.maxSize {
		return rl.rotate()
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
		rl.currentFile = nil
	}

	for i := rl.backupCount - 1; i >= 0; i-- {
		srcPath := rl.getBackupPath(i)
		dstPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(srcPath); err == nil {
			if i == rl.backupCount-1 {
				os.Remove(srcPath)
			} else {
				os.Rename(srcPath, dstPath)
			}
		}
	}

	if _, err := os.Stat(rl.basePath); err == nil {
		backupPath := rl.getBackupPath(0)
		os.Rename(rl.basePath, backupPath)

		if rl.compressOld {
			go rl.compressFile(backupPath)
		}
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d", rl.basePath, index+1)
}

func (rl *RotatingLogger) compressFile(path string) {
	compressedPath := path + ".gz"
	// Compression implementation would go here
	// For now just rename to simulate compression
	if err := os.Rename(path, compressedPath); err != nil {
		log.Printf("Failed to compress %s: %v", path, err)
	}
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}