
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
	mu          sync.Mutex
}

func NewLogRotator(filePath string, maxSize int64) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
	}

	if err := rotator.openFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openFile() error {
	dir := filepath.Dir(lr.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
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

	lr.file = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) rotate() error {
	if lr.file != nil {
		lr.file.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	return lr.openFile()
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

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 10*1024*1024)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}