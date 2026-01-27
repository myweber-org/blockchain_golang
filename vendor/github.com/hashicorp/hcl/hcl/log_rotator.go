
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
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
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

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

func (lr *LogRotator) rotateIfNeeded() error {
	if lr.currentSize < lr.maxSize {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, archivePath); err != nil {
		return err
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) Write(data []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if err := lr.rotateIfNeeded(); err != nil {
		return 0, err
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
	rotator, err := NewLogRotator("application.log", 1024*1024) // 1MB max size
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}