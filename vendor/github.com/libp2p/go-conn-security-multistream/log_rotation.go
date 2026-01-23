package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	maxFiles     int
	currentSize  int64
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
	if err := os.MkdirAll(filepath.Dir(basePath), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		file:        file,
		basePath:    basePath,
		maxSize:     maxSize,
		maxFiles:    maxFiles,
		currentSize: info.Size(),
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	for i := rl.maxFiles - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", rl.basePath, i)
		newPath := fmt.Sprintf("%s.%d", rl.basePath, i+1)
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	backupPath := fmt.Sprintf("%s.1", rl.basePath)
	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}