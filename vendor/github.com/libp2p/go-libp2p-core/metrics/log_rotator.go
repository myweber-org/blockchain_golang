
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

const (
	maxFileSize   = 10 * 1024 * 1024 // 10MB
	backupCount   = 5
	logDir        = "./logs"
	currentLog    = "app.log"
	compressOld   = true
)

type LogRotator struct {
	mu          sync.Mutex
	currentSize int64
	file        *os.File
}

func NewLogRotator() (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(logDir, currentLog)
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
		file:        file,
		currentSize: info.Size(),
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	n, err := lr.file.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	if lr.currentSize >= maxFileSize {
		if err := lr.rotate(); err != nil {
			log.Printf("Failed to rotate log: %v", err)
		}
	}

	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return err
	}

	oldPath := filepath.Join(logDir, currentLog)
	timestamp := time.Now().Format("20060102_150405")
	newPath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.file = file
	lr.currentSize = 0

	go lr.manageBackups(newPath)

	return nil
}

func (lr *LogRotator) manageBackups(newPath string) {
	pattern := filepath.Join(logDir, "app_*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > backupCount {
		toRemove := matches[:len(matches)-backupCount]
		for _, file := range toRemove {
			if compressOld && !strings.HasSuffix(file, ".gz") {
				lr.compressFile(file)
			} else {
				os.Remove(file)
			}
		}
	}
}

func (lr *LogRotator) compressFile(path string) {
	src, err := os.Open(path)
	if err != nil {
		return
	}
	defer src.Close()

	dstPath := path + ".gz"
	dst, err := os.Create(dstPath)
	if err != nil {
		return
	}
	defer dst.Close()

	// Simple compression simulation
	// In real implementation, use gzip.NewWriter(dst)
	_, err = io.Copy(dst, src)
	if err != nil {
		os.Remove(dstPath)
		return
	}

	os.Remove(path)
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(rotator)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("X", 10000))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}