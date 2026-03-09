
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
	currentSize  int64
	maxSize      int64
	basePath     string
	currentFile  *os.File
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s.%d", rl.basePath, timestamp, rl.rotationCount)
	rl.rotationCount++

	if err := os.Rename(rl.basePath, rotatedPath); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		rl.mu.Unlock()
		if err := rl.rotate(); err != nil {
			rl.mu.Lock()
			return 0, err
		}
		rl.mu.Lock()
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
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
	logger, err := NewRotatingLogger("app.log", 1024*1024)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write failed: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}